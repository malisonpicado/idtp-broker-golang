package controller

import (
	"idtp/storage"
	"net"
)

func Broadcast(
	varIndex uint32,
	response []byte,
	depman *storage.DependentsManager,
	clients *storage.ClientsList,
	currentEntity net.Conn) {

	cls := clients.GetClients()

	for _, client := range cls {
		if client == currentEntity {
			continue
		}

		client.Write(response)
	}

	deps := depman.GetDependentsOf(varIndex)

	for _, dependent := range deps {
		if dependent == currentEntity {
			continue
		}

		dependent.Write(response)
	}
}
