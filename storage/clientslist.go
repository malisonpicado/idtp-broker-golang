package storage

import (
	"net"
	"sync"
)

type ClientsList struct {
	Clients map[*net.Conn]struct{}
	Mu      sync.Mutex
}

func InitializeClientsList() *ClientsList {
	return &ClientsList{
		Clients: make(map[*net.Conn]struct{}),
	}
}

func (clients *ClientsList) AddClient(clientConn *net.Conn) {
	clients.Mu.Lock()
	defer clients.Mu.Unlock()

	_, exists := clients.Clients[clientConn]

	if exists {
		panic("Client already registered")
	}

	clients.Clients[clientConn] = struct{}{}
}

func (clients *ClientsList) RemoveClient(clientConn *net.Conn) {
	clients.Mu.Lock()
	defer clients.Mu.Unlock()

	delete(clients.Clients, clientConn)
}

func (clients *ClientsList) GetClients() []*net.Conn {
	clients.Mu.Lock()
	defer clients.Mu.Unlock()

	cls := make([]*net.Conn, 0, len(clients.Clients))

	for client := range clients.Clients {
		cls = append(cls, client)
	}

	return cls
}
