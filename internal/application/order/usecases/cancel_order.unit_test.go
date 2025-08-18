//go:build all || unit || usecase

package usecases_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	orderUsecases "github.com/juninhoitabh/clob-go/internal/application/order/usecases"
	"github.com/juninhoitabh/clob-go/internal/application/order/usecases/fakers"
	domainAccount "github.com/juninhoitabh/clob-go/internal/domain/account"
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	domainOrder "github.com/juninhoitabh/clob-go/internal/domain/order"
	accountMocks "github.com/juninhoitabh/clob-go/internal/infra/repositories/account/mocks"
	bookMocks "github.com/juninhoitabh/clob-go/internal/infra/repositories/book/mocks"
	orderMocks "github.com/juninhoitabh/clob-go/internal/infra/repositories/order/mocks"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type CancelOrderUseCaseUnitTestSuite struct {
	suite.Suite
	inputFaker  orderUsecases.CancelOrderInput
	bookRepo    *bookMocks.MockIBookRepository
	orderRepo   *orderMocks.MockIOrderRepository
	accountRepo *accountMocks.MockIAccountRepository
	ctrl        *gomock.Controller
	usecase     *orderUsecases.CancelOrderUseCase
}

func (suite *CancelOrderUseCaseUnitTestSuite) SetupTest() {
	suite.inputFaker = fakers.CancelOrderInputFaker()
	suite.ctrl = gomock.NewController(suite.T())
	suite.bookRepo = bookMocks.NewMockIBookRepository(suite.ctrl)
	suite.orderRepo = orderMocks.NewMockIOrderRepository(suite.ctrl)
	suite.accountRepo = accountMocks.NewMockIAccountRepository(suite.ctrl)
	suite.usecase = orderUsecases.NewCancelOrderUseCase(suite.bookRepo, suite.orderRepo, suite.accountRepo)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_Success() {
	input := suite.inputFaker

	order := &domainOrder.Order{
		AccountID:  "acc123",
		Instrument: "BTC/USDT",
		Side:       domainOrder.Buy,
		Price:      100,
		Remaining:  5,
	}
	order.ID.ID = input.OrderID
	book := &domainBook.Book{Instrument: "BTC/USDT"}
	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 0, Reserved: 500},
		},
	}

	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(order, nil)
	suite.bookRepo.EXPECT().GetBook(order.Instrument).Return(book, nil)
	suite.bookRepo.EXPECT().SaveBook(book).Return(nil)
	suite.accountRepo.EXPECT().Get(order.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(gomock.Any()).Return(nil).AnyTimes()
	suite.orderRepo.EXPECT().SaveOrder(gomock.Any()).Return(nil).AnyTimes()

	out, err := suite.usecase.Execute(input)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), out)
	assert.Equal(suite.T(), order, out.Order)
	assert.Equal(suite.T(), int64(0), order.Remaining)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_OrderNotFound() {
	input := suite.inputFaker
	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(nil, shared.ErrNotFound)

	out, err := suite.usecase.Execute(input)
	assert.ErrorIs(suite.T(), err, shared.ErrNotFound)
	assert.Nil(suite.T(), out)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_BookNotFound() {
	input := suite.inputFaker
	order := &domainOrder.Order{Instrument: "BTC/USDT"}

	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(order, nil)
	suite.bookRepo.EXPECT().GetBook(order.Instrument).Return(nil, shared.ErrNotFound)

	out, err := suite.usecase.Execute(input)
	assert.ErrorIs(suite.T(), err, shared.ErrNotFound)
	assert.Nil(suite.T(), out)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_AccountNotFound() {
	input := suite.inputFaker
	order := &domainOrder.Order{
		AccountID:  "acc123",
		Instrument: "BTC/USDT",
		Side:       domainOrder.Buy,
		Price:      100,
		Remaining:  5,
	}
	order.ID.ID = input.OrderID
	book := &domainBook.Book{Instrument: "BTC/USDT"}

	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(order, nil)
	suite.bookRepo.EXPECT().GetBook(order.Instrument).Return(book, nil)
	suite.bookRepo.EXPECT().SaveBook(book).Return(nil)
	suite.accountRepo.EXPECT().Get(order.AccountID).Return(nil, shared.ErrNotFound)

	out, err := suite.usecase.Execute(input)
	assert.ErrorIs(suite.T(), err, shared.ErrNotFound)
	assert.Nil(suite.T(), out)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_ReleaseReservedError() {
	input := suite.inputFaker
	order := &domainOrder.Order{
		AccountID:  "acc123",
		Instrument: "BTC/USDT",
		Side:       domainOrder.Buy,
		Price:      100,
		Remaining:  5,
	}
	order.ID.ID = input.OrderID
	book := &domainBook.Book{Instrument: "BTC/USDT"}

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 0, Reserved: 0},
		},
	}

	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(order, nil)
	suite.bookRepo.EXPECT().GetBook(order.Instrument).Return(book, nil)
	suite.bookRepo.EXPECT().SaveBook(book).Return(nil)
	suite.accountRepo.EXPECT().Get(order.AccountID).Return(account, nil)

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_SaveAccountError() {
	input := suite.inputFaker
	order := &domainOrder.Order{
		AccountID:  "acc123",
		Instrument: "BTC/USDT",
		Side:       domainOrder.Buy,
		Price:      100,
		Remaining:  5,
	}
	order.ID.ID = input.OrderID
	book := &domainBook.Book{Instrument: "BTC/USDT"}
	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 0, Reserved: 500},
		},
	}

	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(order, nil)
	suite.bookRepo.EXPECT().GetBook(order.Instrument).Return(book, nil)
	suite.bookRepo.EXPECT().SaveBook(book).Return(nil)
	suite.accountRepo.EXPECT().Get(order.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(gomock.Any()).Return(errors.New("save error"))

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_SaveOrderError() {
	input := suite.inputFaker
	order := &domainOrder.Order{
		AccountID:  "acc123",
		Instrument: "BTC/USDT",
		Side:       domainOrder.Buy,
		Price:      100,
		Remaining:  5,
	}
	order.ID.ID = input.OrderID
	book := &domainBook.Book{Instrument: "BTC/USDT"}
	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 0, Reserved: 500},
		},
	}

	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(order, nil)
	suite.bookRepo.EXPECT().GetBook(order.Instrument).Return(book, nil)
	suite.bookRepo.EXPECT().SaveBook(book).Return(nil)
	suite.accountRepo.EXPECT().Get(order.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(gomock.Any()).Return(nil)
	suite.orderRepo.EXPECT().SaveOrder(order).Return(errors.New("save order error"))

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_OrderIsNil() {
	input := suite.inputFaker
	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(nil, nil)

	out, err := suite.usecase.Execute(input)
	assert.ErrorIs(suite.T(), err, shared.ErrNotFound)
	assert.Nil(suite.T(), out)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_BookIsNil() {
	input := suite.inputFaker
	order := &domainOrder.Order{
		AccountID:  "acc123",
		Instrument: "BTC/USDT",
		Side:       domainOrder.Buy,
		Price:      100,
		Remaining:  5,
	}
	order.ID.ID = input.OrderID
	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(order, nil)
	suite.bookRepo.EXPECT().GetBook(order.Instrument).Return(nil, nil)

	out, err := suite.usecase.Execute(input)
	assert.ErrorIs(suite.T(), err, shared.ErrNotFound)
	assert.Nil(suite.T(), out)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_SplitInstrumentError() {
	input := suite.inputFaker
	order := &domainOrder.Order{
		AccountID:  "acc123",
		Instrument: "INVALID_INSTRUMENT",
		Side:       domainOrder.Buy,
		Price:      100,
		Remaining:  5,
	}
	order.ID.ID = input.OrderID
	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(order, nil)
	suite.bookRepo.EXPECT().GetBook(order.Instrument).Return(&domainBook.Book{}, nil)
	suite.bookRepo.EXPECT().SaveBook(gomock.Any()).Return(nil)

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_ReleaseReservedError_SellOrder() {
	input := suite.inputFaker
	order := &domainOrder.Order{
		AccountID:  "acc123",
		Instrument: "BTC/USDT",
		Side:       domainOrder.Sell,
		Price:      100,
		Remaining:  5,
	}
	order.ID.ID = input.OrderID
	book := &domainBook.Book{Instrument: "BTC/USDT"}
	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"BTC": {Available: 0, Reserved: 0},
		},
	}

	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(order, nil)
	suite.bookRepo.EXPECT().GetBook(order.Instrument).Return(book, nil)
	suite.accountRepo.EXPECT().Get(order.AccountID).Return(account, nil)
	suite.bookRepo.EXPECT().SaveBook(gomock.Any()).Return(nil)

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *CancelOrderUseCaseUnitTestSuite) TestExecute_SaveBookError() {
	input := suite.inputFaker
	order := &domainOrder.Order{
		AccountID:  "acc123",
		Instrument: "BTC/USDT",
		Side:       domainOrder.Buy,
		Price:      100,
		Remaining:  5,
	}
	order.ID.ID = input.OrderID
	book := &domainBook.Book{Instrument: "BTC/USDT"}

	suite.orderRepo.EXPECT().GetOrder(input.OrderID).Return(order, nil)
	suite.bookRepo.EXPECT().GetBook(order.Instrument).Return(book, nil)
	suite.bookRepo.EXPECT().SaveBook(book).Return(errors.New("save book error"))

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CancelOrderUseCaseUnitTestSuite))
	suite.Run(t, new(PlaceOrderUseCaseUnitTestSuite))
}
