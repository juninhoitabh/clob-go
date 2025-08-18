//go:build all || unit || domain

package services_test

import (
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/domain/account/fakers"
	"github.com/juninhoitabh/clob-go/internal/domain/account/services"
	"github.com/juninhoitabh/clob-go/internal/infra/repositories/account/mocks"
	"github.com/juninhoitabh/clob-go/internal/shared"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

type SettleTradeUnitTestSuite struct {
	suite.Suite
	params          fakers.SettleTradeParams
	accountRepoMock *mocks.MockIAccountRepository
	ctrl            *gomock.Controller
	buyer           *account.Account
	seller          *account.Account
}

func (suite *SettleTradeUnitTestSuite) SetupTest() {
	suite.params = fakers.SettleTradeParamsFaker()
	suite.ctrl = gomock.NewController(suite.T())
	suite.accountRepoMock = mocks.NewMockIAccountRepository(suite.ctrl)

	buyerProps := account.AccountProps{Name: "buyer"}
	sellerProps := account.AccountProps{Name: "seller"}
	buyer, _ := account.NewAccount(buyerProps, idObjValue.Uuid)
	seller, _ := account.NewAccount(sellerProps, idObjValue.Uuid)
	buyer.ID.ID = suite.params.BuyerID
	seller.ID.ID = suite.params.SellerID

	buyer.Credit(suite.params.Quote, shared.Mul(suite.params.Price, suite.params.Qty))
	seller.Credit(suite.params.Base, suite.params.Qty)
	buyer.Reserve(suite.params.Quote, shared.Mul(suite.params.Price, suite.params.Qty))
	seller.Reserve(suite.params.Base, suite.params.Qty)

	suite.buyer = buyer
	suite.seller = seller
}

func (suite *SettleTradeUnitTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *SettleTradeUnitTestSuite) TestSettleTrade_Success() {
	params := suite.params

	// Mock Get para buyer e seller
	suite.accountRepoMock.EXPECT().Get(params.BuyerID).Return(suite.buyer, nil)
	suite.accountRepoMock.EXPECT().Get(params.SellerID).Return(suite.seller, nil)
	// Mock Save para buyer e seller
	suite.accountRepoMock.EXPECT().Save(suite.buyer).Return(nil)
	suite.accountRepoMock.EXPECT().Save(suite.seller).Return(nil)

	err := services.SettleTrade(
		suite.accountRepoMock,
		params.BuyerID,
		params.SellerID,
		params.Base,
		params.Quote,
		params.Price,
		params.Qty,
	)
	suite.NoError(err)
}

func (suite *SettleTradeUnitTestSuite) TestSettleTrade_BuyerNotFound() {
	params := suite.params

	suite.accountRepoMock.EXPECT().Get(params.BuyerID).Return(nil, shared.ErrNotFound)

	err := services.SettleTrade(
		suite.accountRepoMock,
		params.BuyerID,
		params.SellerID,
		params.Base,
		params.Quote,
		params.Price,
		params.Qty,
	)
	suite.ErrorIs(err, shared.ErrNotFound)
}

func (suite *SettleTradeUnitTestSuite) TestSettleTrade_SellerNotFound() {
	params := suite.params

	suite.accountRepoMock.EXPECT().Get(params.BuyerID).Return(suite.buyer, nil)
	suite.accountRepoMock.EXPECT().Get(params.SellerID).Return(nil, shared.ErrNotFound)

	err := services.SettleTrade(
		suite.accountRepoMock,
		params.BuyerID,
		params.SellerID,
		params.Base,
		params.Quote,
		params.Price,
		params.Qty,
	)
	suite.ErrorIs(err, shared.ErrNotFound)
}

func (suite *SettleTradeUnitTestSuite) TestSettleTrade_SaveBuyerError() {
	params := suite.params

	suite.accountRepoMock.EXPECT().Get(params.BuyerID).Return(suite.buyer, nil)
	suite.accountRepoMock.EXPECT().Get(params.SellerID).Return(suite.seller, nil)
	suite.accountRepoMock.EXPECT().Save(suite.buyer).Return(errors.New("save buyer error"))

	err := services.SettleTrade(
		suite.accountRepoMock,
		params.BuyerID,
		params.SellerID,
		params.Base,
		params.Quote,
		params.Price,
		params.Qty,
	)
	suite.ErrorContains(err, "save buyer error")
}

func (suite *SettleTradeUnitTestSuite) TestSettleTrade_SaveSellerError() {
	params := suite.params

	suite.accountRepoMock.EXPECT().Get(params.BuyerID).Return(suite.buyer, nil)
	suite.accountRepoMock.EXPECT().Get(params.SellerID).Return(suite.seller, nil)
	suite.accountRepoMock.EXPECT().Save(suite.buyer).Return(nil)
	suite.accountRepoMock.EXPECT().Save(suite.seller).Return(errors.New("save seller error"))

	err := services.SettleTrade(
		suite.accountRepoMock,
		params.BuyerID,
		params.SellerID,
		params.Base,
		params.Quote,
		params.Price,
		params.Qty,
	)
	suite.ErrorContains(err, "save seller error")
}

func (suite *SettleTradeUnitTestSuite) TestSettleTrade_ErrorBuyerUseReserved() {
	params := suite.params
	suite.buyer.Balances[params.Quote] = &account.Balance{Available: 0, Reserved: 0}

	suite.accountRepoMock.EXPECT().Get(params.BuyerID).Return(suite.buyer, nil)
	suite.accountRepoMock.EXPECT().Get(params.SellerID).Return(suite.seller, nil)

	err := services.SettleTrade(
		suite.accountRepoMock,
		params.BuyerID,
		params.SellerID,
		params.Base,
		params.Quote,
		params.Price,
		params.Qty,
	)
	suite.Error(err)
	suite.Contains(err.Error(), "buyer use reserved")
}

type buyerWithCreditError struct {
	*account.Account
}

func (a *buyerWithCreditError) Credit(asset string, amount int64) error {
	return account.ErrInvalidParam
}

func (suite *SettleTradeUnitTestSuite) TestSettleTrade_ErrorSellerUseReserved() {
	params := suite.params
	suite.seller.Balances[params.Base] = &account.Balance{Available: 0, Reserved: 0}

	suite.accountRepoMock.EXPECT().Get(params.BuyerID).Return(suite.buyer, nil)
	suite.accountRepoMock.EXPECT().Get(params.SellerID).Return(suite.seller, nil)
	suite.accountRepoMock.EXPECT().Save(suite.buyer).Return(nil)

	err := services.SettleTrade(
		suite.accountRepoMock,
		params.BuyerID,
		params.SellerID,
		params.Base,
		params.Quote,
		params.Price,
		params.Qty,
	)
	suite.Error(err)
	suite.Contains(err.Error(), "seller use reserved")
}

type sellerWithCreditError struct {
	*account.Account
}

func (a *sellerWithCreditError) Credit(asset string, amount int64) error {
	return account.ErrInvalidParam
}
