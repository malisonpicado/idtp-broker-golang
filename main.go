package main

import (
	"bufio"
	"fmt"
	"idtp/controller"
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
	if configuration.OperationMode == values.OP_MODE_STRICT && len(configuration.Key) == 0 {
		panic("STRICT operation mode requires a key")
	}

	tcpListener, listenerErr := net.Listen("tcp", ":8080")

	if listenerErr != nil {
		panic(fmt.Append([]byte("Error starting tcp server"), listenerErr.Error()))
	}

	defer tcpListener.Close()

	fmt.Println("IDTP Broker. Software version:", softwareVersion)
	fmt.Println("IDTP Protocol Version: v0.5.0-beta")
	fmt.Println("IDTP Protocol Version Byte: 0x00")
	fmt.Println("Server started. Listening:", tcpListener.Addr())

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
		println("Device disconnected")
	}()

	for {
		deadlineErr := (*conn).SetReadDeadline(time.Now().Add(idleTimeout))
		if deadlineErr != nil {
			panic(fmt.Append([]byte("Error while setting up connection's keep alive:"), deadlineErr.Error()))
		}

		n, err := reader.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return
			}

			panic(fmt.Append([]byte("Error while trying to read buffer: "), err.Error()))
		}

		data := buffer[:n]
		processDataAs := controller.ProcessDataAs(data, entity == nil, &configuration)

		if processDataAs.Type == values.PT_REQUEST {
			if entity == nil {
				(*conn).Write(controller.RequestProcessor(data[1:], configuration, conn, entity, database, clients, dependents))
				return
			}

			(*conn).Write(controller.RequestProcessor(data, configuration, conn, entity, database, clients, dependents))
		}

		if processDataAs.Type == values.PT_ERROR {
			(*conn).Write([]byte{processDataAs.PreRcCode})
			return
		}

		if processDataAs.Type == values.PT_PING {
			(*conn).Write([]byte{values.RC_SUCCESS})
		}

		if processDataAs.Type == values.PT_CREATE_CONNECTION {
			connreq, crcode := controller.ConnectionRequestProcessor(data[1:], configuration, database)

			if !(crcode == values.RC_SUCCESS || crcode == values.RC_SUCCESSFUL_CONN_WITH_PREDEFINED_PARAMS) {
				(*conn).Write([]byte{crcode})
				return
			}

			entity = controller.CreateEntityConfig(connreq, configuration.OperationMode == values.OP_MODE_STRICT || crcode == values.RC_SUCCESSFUL_CONN_WITH_PREDEFINED_PARAMS)
			idleTimeout = time.Duration(connreq.KeepAlive) * time.Second

			if entity.EntityType == values.ENTITY_DEVICE {
				controller.SetupDependencies(conn, entity, dependents)
			} else {
				clients.AddClient(conn)
			}

			(*conn).Write([]byte{values.RC_SUCCESS})
		}

		if processDataAs.Type == values.PT_DISCONNECTION {
			return
		}
	}
}
