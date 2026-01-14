package controller

import (
	"idtp/storage"
	"idtp/values"
	"net"
)

func Broadcast(
	request values.Request,
	response []byte,
	dependents *storage.DependentsManager,
	clients *storage.ClientsList,
	currentEntity *net.Conn) {

	cls := clients.GetClients()

	for _, client := range cls {
		if *client == *currentEntity {
			continue
		}

		(*client).Write(response)
	}

	deps := dependents.GetDependentsOf(request.Index)

	for _, dependent := range deps {
		if *dependent == *currentEntity {
			continue
		}

		(*dependent).Write(response)
	}
}
