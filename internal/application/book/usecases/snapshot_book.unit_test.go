package usecases_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	bookUsecases "github.com/juninhoitabh/clob-go/internal/application/book/usecases"
	"github.com/juninhoitabh/clob-go/internal/application/book/usecases/fakers"
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	domainOrder "github.com/juninhoitabh/clob-go/internal/domain/order"
	"github.com/juninhoitabh/clob-go/internal/infra/repositories/book/mocks"
	"github.com/juninhoitabh/clob-go/internal/shared"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

type SnapshotBookUseCaseUnitTestSuite struct {
	suite.Suite
	inputFaker bookUsecases.SnapshotBookInput
	bookRepo   *mocks.MockIBookRepository
	ctrl       *gomock.Controller
	usecase    *bookUsecases.SnapshotBookUseCase
}

func (suite *SnapshotBookUseCaseUnitTestSuite) SetupTest() {
	suite.inputFaker = fakers.SnapshotBookInputFaker()
	suite.ctrl = gomock.NewController(suite.T())
	suite.bookRepo = mocks.NewMockIBookRepository(suite.ctrl)
	suite.usecase = bookUsecases.NewSnapshotBookUseCase(suite.bookRepo)
}

func (suite *SnapshotBookUseCaseUnitTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *SnapshotBookUseCaseUnitTestSuite) TestExecute_Success() {
	input := suite.inputFaker

	mockBook := &domainBook.Book{
		Instrument: input.Instrument,
	}
	mockBook.Prepare(idObjValue.Uuid)
	mockBook.AddOrder(&domainOrder.Order{Side: domainOrder.Buy, Price: 100, Remaining: 5})
	mockBook.AddOrder(&domainOrder.Order{Side: domainOrder.Buy, Price: 100, Remaining: 3})
	mockBook.AddOrder(&domainOrder.Order{Side: domainOrder.Buy, Price: 99, Remaining: 2})
	mockBook.AddOrder(&domainOrder.Order{Side: domainOrder.Sell, Price: 101, Remaining: 4})
	mockBook.AddOrder(&domainOrder.Order{Side: domainOrder.Sell, Price: 102, Remaining: 1})

	suite.bookRepo.EXPECT().GetBook(input.Instrument).Return(mockBook, nil)

	out, err := suite.usecase.Execute(input)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), out)
	assert.Equal(suite.T(), input.Instrument, out.Instrument)
	assert.Len(suite.T(), out.Bids, 2)
	assert.Len(suite.T(), out.Asks, 2)
}

func (suite *SnapshotBookUseCaseUnitTestSuite) TestExecute_BookNotFound() {
	input := suite.inputFaker
	suite.bookRepo.EXPECT().GetBook(input.Instrument).Return(nil, nil)

	out, err := suite.usecase.Execute(input)
	assert.ErrorIs(suite.T(), err, shared.ErrNotFound)
	assert.Nil(suite.T(), out)
}

func (suite *SnapshotBookUseCaseUnitTestSuite) TestExecute_RepoError() {
	input := suite.inputFaker
	suite.bookRepo.EXPECT().GetBook(input.Instrument).Return(nil, errors.New("repo error"))

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SnapshotBookUseCaseUnitTestSuite))
}
