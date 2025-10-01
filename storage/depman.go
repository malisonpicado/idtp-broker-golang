// Variable Dependant Manager
package storage

import "errors"

func (storage *Storage) AddDependents(indexes []uint32, entityId uint32) uint32 {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	errors := uint32(0)
	for _, index := range indexes {
		if index >= storage.Length {
			errors++
			continue
		}

		storage.Data[index].Dependents[entityId] = struct{}{}
	}

	return errors
}

func (storage *Storage) RemoveDependents(indexes []uint32, entityId uint32) uint32 {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	errors := uint32(0)

	for _, index := range indexes {
		if index >= storage.Length {
			errors++
			continue
		}
		delete(storage.Data[index].Dependents, entityId)
	}

	return errors
}

func (storage *Storage) GetDependentsAt(index uint32) ([]uint32, error) {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	if index >= storage.Length {
		return nil, errors.New("invalid index")
	}

	keys := make([]uint32, 0, len(storage.Data[index].Dependents))
	for k := range storage.Data[index].Dependents {
		keys = append(keys, k)
	}

	return keys, nil
}
