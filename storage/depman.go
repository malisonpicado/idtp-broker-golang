// Variable Dependent Manager
package storage

import (
	"net"
	"sync"
)

type dependents = map[*net.Conn]struct{}

type DependentsManager struct {
	Variables map[uint32]dependents // uint32 = index of a variable; dependents = conns that depends of this variable
	Mu        sync.Mutex
}

func InitializeDependencyManager() *DependentsManager {
	return &DependentsManager{
		Variables: make(map[uint32]dependents),
	}
}

func (depman *DependentsManager) AddDependentTo(variableIndex uint32, deviceConn *net.Conn) {
	depman.Mu.Lock()
	defer depman.Mu.Unlock()

	_, exists := depman.Variables[variableIndex]

	if !exists {
		deps := make(dependents)
		deps[deviceConn] = struct{}{}
		depman.Variables[variableIndex] = deps
		return
	}

	var dependents dependents = depman.Variables[variableIndex]
	dependents[deviceConn] = struct{}{}
	depman.Variables[variableIndex] = dependents
}

func (depman *DependentsManager) RemoveDependentFrom(variableIndex uint32, deviceConn *net.Conn) {
	depman.Mu.Lock()
	defer depman.Mu.Unlock()

	_, exists := depman.Variables[variableIndex]

	if !exists {
		return
	}

	delete(depman.Variables[variableIndex], deviceConn)

	if len(depman.Variables[variableIndex]) == 0 {
		delete(depman.Variables, variableIndex)
	}
}

func (depman *DependentsManager) GetDependentsOf(variableIndex uint32) []*net.Conn {
	depman.Mu.Lock()
	defer depman.Mu.Unlock()

	_, exists := depman.Variables[variableIndex]

	if !exists {
		return nil
	}

	var depset dependents = depman.Variables[variableIndex]
	dependents := make([]*net.Conn, 0, len(depset))

	for item := range depset {
		dependents = append(dependents, item)
	}

	return dependents
}
