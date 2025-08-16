package entities_test

import (
	"reflect"
	"testing"

	domain "github.com/juninhoitabh/clob-go/internal/shared/domain/entities"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BaseEntityTestSuite struct {
	suite.Suite
}

func (suite *BaseEntityTestSuite) SetupTest() {
}

func (suite *BaseEntityTestSuite) TestCheckFieldsAndNumberOfStruct() {
	v := new(domain.BaseEntity)
	metaValue := reflect.ValueOf(v).Elem()

	structFieldsToCheck := []string{"ID"}
	numberStructFields := metaValue.NumField()

	for _, name := range structFieldsToCheck {
		field := metaValue.FieldByName(name)
		if !field.IsValid() {
			suite.T().Errorf("Field %s not exist in struct", name)
		}
	}

	require.Equal(suite.T(), numberStructFields, len(structFieldsToCheck))
}

func (suite *BaseEntityTestSuite) TestNewBaseEntityWithDefaultIDTypeObjectID() {
	baseEntity := domain.BaseEntity{}
	baseEntity.NewBaseEntity("", idObjValue.ObjectID)

	require.NotNil(suite.T(), baseEntity)
	require.IsType(suite.T(), "string", baseEntity.ID.ID)
}

func (suite *BaseEntityTestSuite) TestNewBaseEntityWithSetID() {
	id := idObjValue.NewID("", idObjValue.ObjectID)

	baseEntity := domain.BaseEntity{}
	baseEntity.NewBaseEntity(id.ID, idObjValue.ObjectID)

	require.NotNil(suite.T(), baseEntity)
	require.IsType(suite.T(), "string", baseEntity.ID.ID)
	require.Equal(suite.T(), id.ID, baseEntity.ID.ID)
}

func (suite *BaseEntityTestSuite) TestNewBaseEntityWithSetIDAndNotReset() {
	id := idObjValue.NewID("", idObjValue.ObjectID)

	baseEntity := domain.BaseEntity{}
	baseEntity.NewBaseEntity(id.ID, idObjValue.ObjectID)
	baseEntity.NewBaseEntity("", idObjValue.ObjectID)

	require.NotNil(suite.T(), baseEntity)
	require.IsType(suite.T(), "string", baseEntity.ID.ID)
	require.Equal(suite.T(), id.ID, baseEntity.ID.ID)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(BaseEntityTestSuite))
}
