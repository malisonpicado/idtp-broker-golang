package controller

import (
	"idtp/storage"
	"idtp/values"
	"net"
)

func SetupDependencies(conn *net.Conn, entity *values.Entity, depman *storage.DependentsManager) {
	for _, dep := range entity.DependencyParams {
		depman.AddDependentTo(dep, conn)
	}
}

func ClearDependencies(conn *net.Conn, entity *values.Entity, depman *storage.DependentsManager) {
	for _, dep := range entity.DependencyParams {
		depman.RemoveDependentFrom(dep, conn)
	}
}

func CreateEntityConfig(connreq values.ConnectionRequest, setStrict bool) *values.Entity {
	entity := &values.Entity{
		EntityType: connreq.EntityType,
	}

	if connreq.EntityType == values.ENTITY_CLIENT {
		entity.ProcessAsStrict = false
		return entity
	}

	for _, param := range connreq.Parameters {
		if param.Method == 0 {
			entity.DependencyParams = append(entity.DependencyParams, param.Index)
			continue
		}

		entity.UpdateParams = append(entity.UpdateParams, param.Index)
	}

	entity.ProcessAsStrict = setStrict
	return entity
}
