package storage

import (
	"errors"
	"idtp/values"
	"net"
	"sync"
	"time"
)

type ClientsList struct {
	ClientsIDs map[uint32]struct{}
}

type Entity struct {
	Connection   *net.Conn
	Id           uint32
	EntityType   values.EntityType
	KeepAlive    time.Duration
	IsClosed     bool
	Updating     []uint32
	Dependencies []uint32
}

type EntitiesList struct {
	Mu       sync.RWMutex
	Entities map[uint32]*Entity
}

func (list *EntitiesList) AddEntity(entity *Entity) error {
	list.Mu.Lock()
	defer list.Mu.Unlock()

	_, exist := list.Entities[entity.Id]

	if exist {
		return errors.New("entity already exists")
	}

	list.Entities[entity.Id] = entity
	return nil
}

func (list *EntitiesList) RemoveEntity(entity *Entity) {
	list.Mu.Lock()
	defer list.Mu.Unlock()

	delete(list.Entities, entity.Id)
}

func (list *ClientsList) AddClient(clientId uint32) error {
	_, exists := list.ClientsIDs[clientId]

	if exists {
		return errors.New("client already exists")
	}

	list.ClientsIDs[clientId] = struct{}{}
	return nil
}

func (list *ClientsList) RemoveClient(clientId uint32) {
	delete(list.ClientsIDs, clientId)
}
