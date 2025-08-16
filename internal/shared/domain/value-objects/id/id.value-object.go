package id

import (
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TypeIdEnum string

const (
	ObjectID TypeIdEnum = "ObjectID"
	Uuid     TypeIdEnum = "Uuid"
	Str      TypeIdEnum = "String"
)

type ID struct {
	ID string
}

func NewID(id string, typeId TypeIdEnum) ID {
	if typeId == ObjectID {
		if id == "" {
			return ID{ID: primitive.NewObjectID().Hex()}
		}

		if newObjectID, err := primitive.ObjectIDFromHex(id); err == nil {
			return ID{ID: newObjectID.Hex()}
		}

		return ID{ID: primitive.NewObjectID().Hex()}
	}

	if typeId == Uuid {
		if id == "" {
			return ID{ID: uuid.NewV4().String()}
		}

		if newUuid, err := uuid.FromString(id); err == nil {
			return ID{ID: newUuid.String()}
		}

		return ID{ID: uuid.NewV4().String()}
	}

	return ID{
		ID: id,
	}
}
