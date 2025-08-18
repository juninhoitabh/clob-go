//go:build all || unit || domain

package services_test

import (
	"testing"

	"github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/domain/account/services"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
	"github.com/stretchr/testify/suite"
)

type TransferUnitTestSuite struct {
	suite.Suite
	from *account.Account
	to   *account.Account
}

func (suite *TransferUnitTestSuite) SetupTest() {
	fromProps := account.AccountProps{Name: "from"}
	toProps := account.AccountProps{Name: "to"}
	from, _ := account.NewAccount(fromProps, idObjValue.Uuid)
	to, _ := account.NewAccount(toProps, idObjValue.Uuid)

	suite.from = from
	suite.to = to
}

func (suite *TransferUnitTestSuite) TestTransfer_Success() {
	asset := "BTC"
	amount := int64(100)
	suite.from.Credit(asset, amount)

	err := services.Transfer(suite.from, suite.to, asset, amount)
	suite.NoError(err)
	suite.Equal(int64(0), suite.from.Balances[asset].Available)
	suite.Equal(int64(100), suite.to.Balances[asset].Available)
}

func (suite *TransferUnitTestSuite) TestTransfer_InsufficientFunds() {
	asset := "BTC"
	amount := int64(50) // from não tem saldo

	err := services.Transfer(suite.from, suite.to, asset, amount)
	suite.ErrorIs(err, account.ErrInsufficient)
}

func (suite *TransferUnitTestSuite) TestTransfer_InvalidAmount() {
	asset := "BTC"
	amount := int64(-10)
	suite.from.Credit(asset, 100)

	err := services.Transfer(suite.from, suite.to, asset, amount)
	suite.ErrorIs(err, account.ErrInvalidParam)
}

func (suite *TransferUnitTestSuite) TestTransfer_UseReservedError() {
	asset := "BTC"
	amount := int64(100)
	suite.from.Credit(asset, amount)
	suite.from.Reserve(asset, 50)

	err := services.Transfer(suite.from, suite.to, asset, amount)
	suite.ErrorIs(err, account.ErrInsufficient)
}

func (suite *TransferUnitTestSuite) TestTransfer_CreditError() {
	asset := "BTC"
	amount := int64(-100)
	suite.from.Credit(asset, amount)

	err := services.Transfer(suite.from, suite.to, asset, amount)
	suite.ErrorIs(err, account.ErrInvalidParam)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SettleTradeUnitTestSuite))
	suite.Run(t, new(TransferUnitTestSuite))
}
