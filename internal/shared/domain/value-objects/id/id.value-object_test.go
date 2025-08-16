package valueObjects_test

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
	v := new(idObjValue.ID)
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

func (suite *ValueObjectIDTestSuite) TestShouldCreateAnDefaultObjectID() {
	id := idObjValue.NewID("", idObjValue.ObjectID)

	require.NotNil(suite.T(), id)
	require.IsType(suite.T(), "string", id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldCreateAnObjectIDPassedAnInvalidIdAsParam() {
	newId := fmt.Sprint(suite.faker.Number(0, 9))
	id := idObjValue.NewID(newId, idObjValue.ObjectID)

	require.NotNil(suite.T(), id)
	require.NotEqual(suite.T(), newId, id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldReturnTheObjectIDPassedIdAsParam() {
	newId := primitive.NewObjectID().Hex()
	id := idObjValue.NewID(newId, idObjValue.ObjectID)

	require.NotNil(suite.T(), id)
	require.Equal(suite.T(), newId, id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldCreateAnDefaultUUID() {
	id := idObjValue.NewID("", idObjValue.Uuid)

	require.NotNil(suite.T(), id)
	require.IsType(suite.T(), "string", id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldCreateAnUUIDPassedAnInvalidIdAsParam() {
	newId := fmt.Sprint(suite.faker.Number(0, 9))
	id := idObjValue.NewID(newId, idObjValue.Uuid)

	require.NotNil(suite.T(), id)
	require.NotEqual(suite.T(), newId, id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldReturnTheUUIDPassedIdAsParam() {
	newId := uuid.NewV4().String()
	id := idObjValue.NewID(newId, idObjValue.Uuid)

	require.NotNil(suite.T(), id)
	require.Equal(suite.T(), newId, id.ID)
}

func (suite *ValueObjectIDTestSuite) TestShouldCreateAnIDPassedAnValidIdStringAsParam() {
	newId := fmt.Sprint(suite.faker.Number(0, 9))
	id := idObjValue.NewID(newId, idObjValue.Str)

	require.NotNil(suite.T(), id)
	require.Equal(suite.T(), newId, id.ID)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ValueObjectIDTestSuite))
}
