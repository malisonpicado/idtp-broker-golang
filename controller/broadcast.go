package controller

import (
	"idtp/parsers"
	"idtp/storage"
)

// Broadcast envía a todos los clientes y a los dispositivos dependientes
// las actualizaciones de las variables indicadas en `requests`.
func Broadcast(requests []parsers.Request, senderId uint32, omitSender bool,
	entities *storage.EntitiesList,
	clients *storage.ClientsList,
	db *storage.Storage) {
	if len(requests) == 0 {
		return
	}

	// Mapa para evitar duplicados
	targets := make(map[uint32]*storage.Entity)

	entities.Mu.RLock()
	// 1. Agregar todos los clientes
	for clientID := range clients.ClientsIDs {
		// Si omitSender o coincide con senderID, no incluir
		if omitSender || clientID == senderId {
			continue
		}
		if ent, ok := entities.Entities[clientID]; ok && !ent.IsClosed {
			targets[clientID] = ent
		}
	}
	entities.Mu.RUnlock()

	// 2. Recorrer requests para agregar dependientes
	for _, req := range requests {
		dependents, err := db.GetDependentsAt(req.Index)
		if err != nil {
			continue // si hay error, omitimos dependientes de esta variable
		}
		entities.Mu.RLock()
		for _, depID := range dependents {
			if omitSender || depID == senderId {
				continue
			}
			if ent, ok := entities.Entities[depID]; ok && !ent.IsClosed {
				targets[depID] = ent
			}
		}
		entities.Mu.RUnlock()
	}

	// 3. Construir streams una sola vez por request
	streams := make([][]byte, len(requests))
	for i, req := range requests {
		streams[i] = parsers.BuildUpdateStream(byte(req.Type), req.Index, req.Payload)
	}

	// 4. Enviar a cada entidad
	for _, ent := range targets {
		if ent.Connection == nil {
			continue
		}
		for _, data := range streams {
			(*ent.Connection).Write(data)
		}
	}
}
