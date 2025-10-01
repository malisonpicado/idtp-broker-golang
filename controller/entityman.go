package controller

import (
	"idtp/parsers"
	"idtp/storage"
	"idtp/values"
	"net"
	"sync/atomic"
	"time"
)

func CreateEntity(connection *net.Conn, IdCounter *uint32, connReq parsers.ConnectionRequest) *storage.Entity {
	return &storage.Entity{
		Connection:   connection,
		Id:           atomic.AddUint32(IdCounter, 1),
		EntityType:   connReq.EntityType,
		KeepAlive:    time.Duration(connReq.KeepAlive),
		IsClosed:     false,
		Updating:     connReq.ParamUpdates,
		Dependencies: connReq.ParamDepends,
	}
}

// Sets up a new entity. This entity can be a client or a device.
// if client, the adds it to the clients list; if device adds
// all device dependencies to variable. Returns error if any
// dependency has an invalid index.
func SetupEntity(entity *storage.Entity, entitiesList *storage.EntitiesList, clientsList *storage.ClientsList, database *storage.Storage) values.StatusCode {
	if entity.EntityType == values.ENTITY_CLIENT {
		clientsList.AddClient(entity.Id)
		e := entitiesList.AddEntity(entity)

		if e != nil {
			panic(e)
		}

		return values.SUCCESS
	}

	err := database.AddDependents(entity.Dependencies, entity.Id)

	if err != 0 {
		return values.INVALID_INDEX
	}

	e := entitiesList.AddEntity(entity)

	if e != nil {
		panic(e)
	}

	return values.SUCCESS
}

// Close connection after using this function
func CleanEntity(entity *storage.Entity, entitiesList *storage.EntitiesList, clientsList *storage.ClientsList, database *storage.Storage) {
	if entity == nil {
		return
	}

	entitiesList.RemoveEntity(entity)

	if entity.EntityType == values.ENTITY_CLIENT {
		clientsList.RemoveClient(entity.Id)
		return
	}

	database.RemoveDependents(entity.Dependencies, entity.Id)
}
