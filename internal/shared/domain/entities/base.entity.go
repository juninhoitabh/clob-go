package entities

import (
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

type BaseEntity struct {
	idObjValue.ID
}

func (baseEntity *BaseEntity) NewBaseEntity(id string, typeId idObjValue.TypeIdEnum) {
	if baseEntity.ID.ID != "" {
		return
	}

	baseEntity.ID = idObjValue.NewID(id, typeId)
}
