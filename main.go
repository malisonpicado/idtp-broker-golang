package main

import (
	"bufio"
	"idtp/controller"
	"idtp/storage"
	"idtp/values"
	"net"
	"sync"
	"time"
)

var softwareVersion string = "v0.0.1"

var configuration = values.Configuration{
	Key:             "password123",
	ProtocolVersion: 0xFF,
	OperationMode:   values.OP_MODE_STRICT,
	StorageLimit:    1_000_000,
}

var IdCounter uint32

var database = storage.Initialize(configuration)
var entities = storage.EntitiesList{
	Mu:       sync.RWMutex{},
	Entities: nil,
}
var clients = storage.ClientsList{
	ClientsIDs: map[uint32]struct{}{},
}

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
	connection := *conn
	isFirstConn := true
	idleTimeout := time.Duration(6 * time.Second)
	buffer := make([]byte, 1024)
	reader := bufio.NewReader(connection)
	var entity *storage.Entity = nil

	defer func() {
		controller.CleanEntity(entity, &entities, &clients, database)
		connection.Close()
	}()

	for {
		deadlineErr := connection.SetReadDeadline(time.Now().Add(idleTimeout))
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
			if data[0] == 0xFF {
				connection.Write(controller.RequestProcessor(data[1:], false, database, configuration, &entities, &clients, 0, true))
				return
			}

			if data[0] != 0x00 {
				connection.Write([]byte{byte(values.UNKNOWN_METHOD)})
				return
			}

			connreq, crerr := controller.ConnectionRequestProcessor(data[1:], configuration)

			if crerr != values.SUCCESS {
				connection.Write([]byte{byte(crerr)})
				return
			}

			isFirstConn = false
			entity = controller.CreateEntity(conn, &IdCounter, connreq)
			idleTimeout = time.Duration(entity.KeepAlive * time.Second)

			controller.SetupEntity(entity, &entities, &clients, database)
			connection.Write([]byte{byte(values.SUCCESS)})
			continue
		}

		if len(data) == 1 {
			if data[0] == 0x01 { // PING
				connection.Write([]byte{byte(values.SUCCESS)})
				continue
			}

			if data[0] == 0x02 {
				return
			}

			connection.Write([]byte{byte(values.UNKNOWN_METHOD)})
			continue
		}

		response := controller.RequestProcessor(data, entity.EntityType == values.ENTITY_CLIENT, database, configuration, &entities, &clients, entity.Id, false)
		connection.Write(response)
	}
}
