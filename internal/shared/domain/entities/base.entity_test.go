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
	t := suite.Suite.T()

	v := new(domain.BaseEntity)
	metaValue := reflect.ValueOf(v).Elem()

	structFieldsToCheck := []string{"ID"}
	numberStructFields := metaValue.NumField()

	for _, name := range structFieldsToCheck {
		field := metaValue.FieldByName(name)
		if !field.IsValid() {
			t.Errorf("Field %s not exist in struct", name)
		}
	}

	require.Equal(t, numberStructFields, len(structFieldsToCheck))
}

func (suite *BaseEntityTestSuite) TestNewBaseEntityWithDefaultIDTypeObjectID() {
	t := suite.Suite.T()

	baseEntity := domain.BaseEntity{}
	baseEntity.NewBaseEntity("", idObjValue.ObjectID)

	require.NotNil(t, baseEntity)
	require.IsType(t, "string", baseEntity.ID.ID)
	require.IsType(t, "string", baseEntity.GetID())
}

func (suite *BaseEntityTestSuite) TestNewBaseEntityWithSetID() {
	t := suite.Suite.T()

	id := idObjValue.NewID("", idObjValue.ObjectID)

	baseEntity := domain.BaseEntity{}
	baseEntity.NewBaseEntity(id.ID, idObjValue.ObjectID)

	require.NotNil(t, baseEntity)
	require.IsType(t, "string", baseEntity.ID.ID)
	require.Equal(t, id.ID, baseEntity.ID.ID)
}

func (suite *BaseEntityTestSuite) TestNewBaseEntityWithSetIDAndNotReset() {
	t := suite.Suite.T()

	id := idObjValue.NewID("", idObjValue.ObjectID)

	baseEntity := domain.BaseEntity{}
	baseEntity.NewBaseEntity(id.ID, idObjValue.ObjectID)
	baseEntity.NewBaseEntity("", idObjValue.ObjectID)

	require.NotNil(t, baseEntity)
	require.IsType(t, "string", baseEntity.ID.ID)
	require.Equal(t, id.ID, baseEntity.ID.ID)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(BaseEntityTestSuite))
}
