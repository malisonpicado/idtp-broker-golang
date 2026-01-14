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
