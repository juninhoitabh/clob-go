//go:build all || unit || usecase

package usecases_test

import (
	"errors"

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

type PlaceOrderUseCaseUnitTestSuite struct {
	suite.Suite
	inputFaker  orderUsecases.PlaceOrderInput
	bookRepo    *bookMocks.MockIBookRepository
	orderRepo   *orderMocks.MockIOrderRepository
	accountRepo *accountMocks.MockIAccountRepository
	ctrl        *gomock.Controller
	usecase     *orderUsecases.PlaceOrderUseCase
}

func (suite *PlaceOrderUseCaseUnitTestSuite) SetupTest() {
	suite.inputFaker = fakers.PlaceOrderInputFaker()
	suite.ctrl = gomock.NewController(suite.T())
	suite.bookRepo = bookMocks.NewMockIBookRepository(suite.ctrl)
	suite.orderRepo = orderMocks.NewMockIOrderRepository(suite.ctrl)
	suite.accountRepo = accountMocks.NewMockIAccountRepository(suite.ctrl)
	suite.usecase = orderUsecases.NewPlaceOrderUseCase(suite.bookRepo, suite.orderRepo, suite.accountRepo)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_Success() {
	input := suite.inputFaker
	input.Side = "buy"
	input.Price = 100
	input.Qty = 10

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 1000, Reserved: 0},
		},
	}

	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(account).Return(nil)
	suite.orderRepo.EXPECT().SaveOrder(gomock.Any()).Return(nil)

	book, _ := domainBook.NewBook(domainBook.BookProps{Instrument: input.Instrument}, "Uuid")
	suite.bookRepo.EXPECT().GetBook(input.Instrument).Return(book, nil)
	suite.bookRepo.EXPECT().SaveBook(book).Return(nil).AnyTimes()

	out, err := suite.usecase.Execute(input)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), out)
	assert.Equal(suite.T(), input.AccountID, out.Order.AccountID)
	assert.Equal(suite.T(), input.Instrument, out.Order.Instrument)
	assert.Equal(suite.T(), input.Price, out.Order.Price)
	assert.Equal(suite.T(), input.Qty, out.Order.Qty)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_InvalidParam() {
	input := suite.inputFaker
	input.Price = 0

	out, err := suite.usecase.Execute(input)
	assert.ErrorIs(suite.T(), err, shared.ErrInvalidParam)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_ParseSideError() {
	input := suite.inputFaker
	input.Side = "invalid"

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_AccountNotFound() {
	input := suite.inputFaker

	suite.accountRepo.EXPECT().Get(input.AccountID).Return(nil, errors.New("not found"))

	out, err := suite.usecase.Execute(input)
	assert.ErrorIs(suite.T(), err, shared.ErrNotFound)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_SplitInstrumentError() {
	input := suite.inputFaker
	input.Instrument = "INVALID"

	account := &domainAccount.Account{}
	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_ReserveError_Buy() {
	input := suite.inputFaker
	input.Side = "buy"
	input.Price = 100
	input.Qty = 10

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 0, Reserved: 0},
		},
	}
	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_ReserveError_Sell() {
	input := suite.inputFaker
	input.Side = "sell"
	input.Price = 100
	input.Qty = 10

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"BTC": {Available: 0, Reserved: 0},
		},
	}
	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_NewOrderError() {
	input := suite.inputFaker
	input.Side = "buy"
	input.Price = 100
	input.Qty = 10

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 1000, Reserved: 0},
		},
	}
	suite.accountRepo.EXPECT().Get(gomock.Any()).Return(account, nil).AnyTimes()
	suite.accountRepo.EXPECT().Save(gomock.Any()).Return(nil).AnyTimes()

	input.Qty = 0

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_GetBookError() {
	input := suite.inputFaker
	input.Side = "buy"
	input.Price = 100
	input.Qty = 10

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 1000, Reserved: 0},
		},
	}
	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(account).Return(nil)
	suite.orderRepo.EXPECT().SaveOrder(gomock.Any()).Return(nil)
	suite.bookRepo.EXPECT().GetBook(input.Instrument).Return(nil, errors.New("get book error"))

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_SaveBookError_NewBook() {
	input := suite.inputFaker
	input.Side = "buy"
	input.Price = 100
	input.Qty = 10

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 1000, Reserved: 0},
		},
	}
	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(account).Return(nil)
	suite.orderRepo.EXPECT().SaveOrder(gomock.Any()).Return(nil)
	suite.bookRepo.EXPECT().GetBook(input.Instrument).Return(nil, nil)
	suite.bookRepo.EXPECT().SaveBook(gomock.Any()).Return(errors.New("save book error"))

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_SaveBookError_AfterMatch() {
	input := suite.inputFaker
	input.Side = "buy"
	input.Price = 100
	input.Qty = 10

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 1000, Reserved: 0},
		},
	}
	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(account).Return(nil)
	suite.orderRepo.EXPECT().SaveOrder(gomock.Any()).Return(nil)
	book, _ := domainBook.NewBook(domainBook.BookProps{Instrument: input.Instrument}, "Uuid")
	suite.bookRepo.EXPECT().GetBook(input.Instrument).Return(book, nil)
	suite.bookRepo.EXPECT().SaveBook(gomock.Any()).Return(errors.New("save book error"))

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_SaveAccountError() {
	input := suite.inputFaker
	input.Side = "buy"
	input.Price = 100
	input.Qty = 10

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 1000, Reserved: 0},
		},
	}
	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(account).Return(errors.New("save account error"))

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_SaveAccountError2() {
	input := suite.inputFaker

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 1000, Reserved: 0},
		},
	}
	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(gomock.Any()).Return(errors.New("save error")).AnyTimes()

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_SaveOrderError() {
	input := suite.inputFaker
	input.Side = "buy"
	input.Price = 100
	input.Qty = 10

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 1000, Reserved: 0},
		},
	}
	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(account).Return(nil)
	suite.orderRepo.EXPECT().SaveOrder(gomock.Any()).Return(errors.New("save order error"))

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}

func (suite *PlaceOrderUseCaseUnitTestSuite) TestExecute_SettleTradeError() {
	input := suite.inputFaker
	input.Side = "buy"
	input.Price = 100
	input.Qty = 10

	account := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"USDT": {Available: 1000, Reserved: 0},
		},
	}
	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(account).Return(nil)
	suite.orderRepo.EXPECT().SaveOrder(gomock.Any()).Return(nil)

	book, _ := domainBook.NewBook(domainBook.BookProps{Instrument: input.Instrument}, "Uuid")
	sellerOrder := &domainOrder.Order{
		AccountID: "seller123",
		Side:      domainOrder.Sell,
		Price:     100,
		Qty:       10,
		Remaining: 10,
	}
	book.AddOrder(sellerOrder)
	suite.bookRepo.EXPECT().GetBook(input.Instrument).Return(book, nil)
	suite.bookRepo.EXPECT().SaveBook(book).Return(nil).AnyTimes()

	sellerAccount := &domainAccount.Account{
		Balances: map[string]*domainAccount.Balance{
			"BTC": {Available: 0, Reserved: 0},
		},
	}

	suite.accountRepo.EXPECT().Get(gomock.Any()).Return(sellerAccount, nil).AnyTimes()

	out, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), out)
}
