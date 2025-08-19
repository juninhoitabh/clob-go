package id_test

import (
	"fmt"
	"reflect"
	"testing"

	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"

	faker "github.com/brianvoe/gofakeit/v7"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ValueObjectIDTestSuite struct {
	suite.Suite
	faker *faker.Faker
}

func (suite *ValueObjectIDTestSuite) SetupTest() {
	suite.faker = faker.New(0)
}

func (suite *ValueObjectIDTestSuite) TestCheckFieldsAndNumberOfStruct() {
	t := suite.Suite.T()

	v := new(idObjValue.ID)
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

func (suite *ValueObjectIDTestSuite) TestShouldCreateAnDefaultObjectID() {
	t := suite.Suite.T()

	id := idObjValue.NewID("", idObjValue.ObjectID)

	require.NotNil(t, id)
	require.IsType(t, "string", id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldCreateAnObjectIDPassedAnInvalidIdAsParam() {
	t := suite.Suite.T()

	newId := fmt.Sprint(suite.faker.Number(0, 9))
	id := idObjValue.NewID(newId, idObjValue.ObjectID)

	require.NotNil(t, id)
	require.NotEqual(t, newId, id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldReturnTheObjectIDPassedIdAsParam() {
	t := suite.Suite.T()

	newId := primitive.NewObjectID().Hex()
	id := idObjValue.NewID(newId, idObjValue.ObjectID)

	require.NotNil(t, id)
	require.Equal(t, newId, id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldCreateAnDefaultUUID() {
	t := suite.Suite.T()

	id := idObjValue.NewID("", idObjValue.Uuid)

	require.NotNil(t, id)
	require.IsType(t, "string", id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldCreateAnUUIDPassedAnInvalidIdAsParam() {
	t := suite.Suite.T()

	newId := fmt.Sprint(suite.faker.Number(0, 9))
	id := idObjValue.NewID(newId, idObjValue.Uuid)

	require.NotNil(t, id)
	require.NotEqual(t, newId, id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldReturnTheUUIDPassedIdAsParam() {
	t := suite.Suite.T()

	newId := uuid.NewV4().String()
	id := idObjValue.NewID(newId, idObjValue.Uuid)

	require.NotNil(t, id)
	require.Equal(t, newId, id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldCreateAnIDPassedAnValidIdStringAsParam() {
	t := suite.Suite.T()

	newId := fmt.Sprint(suite.faker.Number(0, 9))
	id := idObjValue.NewID(newId, idObjValue.Str)

	require.NotNil(t, id)
	require.Equal(t, newId, id.ID)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ValueObjectIDTestSuite))
}
