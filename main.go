package main

//

import (
	"bufio"
	"idtp/controller"
	"idtp/parsers"
	"idtp/storage"
	"idtp/values"
	"net"
	"time"
)

var softwareVersion string = "v0.1.0"

var configuration = values.Configuration{
	Key:             "password123",
	ProtocolVersion: 0x00,
	OperationMode:   values.OP_MODE_STRICT,
	StorageLimit:    1_000_000,
}

var IdCounter uint32

var database = storage.InitializeStorage(configuration.StorageLimit)
var dependents = storage.InitializeDependencyManager()
var clients = storage.InitializeClientsList()

func main() {
	println(softwareVersion)
	tcpListener, listenerErr := net.Listen("tcp", ":8080")

	if listenerErr != nil {
		panic("error starting tcp server")
	}

	defer tcpListener.Close()
	println("tcp server started")

	for {
		connection, connErr := tcpListener.Accept()

		if connErr != nil {
			println("device disconnected")
			continue
		}

		println("new connection accepted")
		go handleTcpConnection(&connection)
	}
}

func handleTcpConnection(conn *net.Conn) {
	isFirstConn := true
	idleTimeout := time.Duration(6 * time.Second)
	var entity *values.Entity
	buffer := make([]byte, 1024)
	reader := bufio.NewReader(*conn)

	defer func() {
		if entity != nil {
			if entity.EntityType == values.ENTITY_CLIENT {
				clients.RemoveClient(conn)
			} else {
				controller.ClearDependencies(conn, entity, dependents)
			}
		}

		(*conn).Close()
	}()

	for {
		deadlineErr := (*conn).SetReadDeadline(time.Now().Add(idleTimeout))
		if deadlineErr != nil {
			return
		}

		n, err := reader.Read(buffer)
		if err != nil {
			// if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// 	return
			// }
			return
		}

		data := buffer[:n]

		if isFirstConn {
			if data[0] == 0xFF && configuration.OperationMode != values.OP_MODE_STRICT {
				(*conn).Write(controller.RequestProcessor(data[1:], false, configuration, conn, nil, database, clients, dependents))
				return
			}

			if data[0] != 0x00 {
				(*conn).Write([]byte{byte(values.RC_UNKNOWN_METHOD)})
				return
			}

			connreq, crerr := controller.ConnectionRequestProcessor(data[1:], configuration)

			if crerr != values.RC_SUCCESS {
				(*conn).Write([]byte{crerr})
				return
			}

			isFirstConn = false
			entity = parsers.EntityConfigBuilder(connreq)
			idleTimeout = time.Duration(connreq.KeepAlive) * time.Second

			if entity.EntityType == values.ENTITY_DEVICE {
				controller.SetupDependencies(conn, entity, dependents)
			} else {
				clients.AddClient(conn)
			}

			(*conn).Write([]byte{byte(values.RC_SUCCESS)})
			continue
		}

		if len(data) == 1 {
			if data[0] == 0x01 { // PING
				(*conn).Write([]byte{byte(values.RC_SUCCESS)})
				continue
			}

			if data[0] == 0x02 { // DISCONNECTION
				return
			}

			(*conn).Write([]byte{byte(values.RC_UNKNOWN_METHOD)})
			continue
		}

		response := controller.RequestProcessor(data, entity.EntityType == values.ENTITY_CLIENT, configuration, conn, entity, database, clients, dependents)
		(*conn).Write(response)
	}
}
