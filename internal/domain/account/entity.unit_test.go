//go:build all || unit || domain

package account_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/domain/account/fakers"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

type AccountTestSuite struct {
	suite.Suite
	propsFaker account.AccountProps
}

func (suite *AccountTestSuite) SetupTest() {
	suite.propsFaker = fakers.AccountPropsFaker()
}

func (suite *AccountTestSuite) TestNewAccount_Success() {
	acc, err := account.NewAccount(suite.propsFaker, idObjValue.Uuid)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), acc)
	assert.Equal(suite.T(), suite.propsFaker.Name, acc.Name)
	assert.NotEmpty(suite.T(), acc.CreatedAt)
	assert.NotNil(suite.T(), acc.Balances)
}

func (suite *AccountTestSuite) TestValidate_ErrorOnEmptyName() {
	props := account.AccountProps{Name: ""}
	acc := account.Account{Name: props.Name}
	err := acc.Validate()
	assert.ErrorIs(suite.T(), err, account.ErrInvalidParam)
}

func (suite *AccountTestSuite) TestCredit_Success() {
	acc, _ := account.NewAccount(suite.propsFaker, idObjValue.Uuid)
	err := acc.Credit("BTC", 100)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(100), acc.Balances["BTC"].Available)
}

func (suite *AccountTestSuite) TestCredit_ErrorOnNegativeAmount() {
	acc, _ := account.NewAccount(suite.propsFaker, idObjValue.Uuid)
	err := acc.Credit("BTC", -10)
	assert.ErrorIs(suite.T(), err, account.ErrInvalidParam)
}

func (suite *AccountTestSuite) TestReserve_Success() {
	acc, _ := account.NewAccount(suite.propsFaker, idObjValue.Uuid)
	_ = acc.Credit("ETH", 50)
	err := acc.Reserve("ETH", 30)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(20), acc.Balances["ETH"].Available)
	assert.Equal(suite.T(), int64(30), acc.Balances["ETH"].Reserved)
}

func (suite *AccountTestSuite) TestReserve_ErrorOnInsufficientFunds() {
	acc, _ := account.NewAccount(suite.propsFaker, idObjValue.Uuid)
	_ = acc.Credit("ETH", 10)
	err := acc.Reserve("ETH", 20)
	assert.ErrorIs(suite.T(), err, account.ErrInsufficient)
}

func (suite *AccountTestSuite) TestUseReserved_Success() {
	acc, _ := account.NewAccount(suite.propsFaker, idObjValue.Uuid)
	_ = acc.Credit("USDT", 100)
	_ = acc.Reserve("USDT", 60)
	err := acc.UseReserved("USDT", 50)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(10), acc.Balances["USDT"].Reserved)
}

func (suite *AccountTestSuite) TestUseReserved_ErrorOnInsufficientReserved() {
	acc, _ := account.NewAccount(suite.propsFaker, idObjValue.Uuid)
	_ = acc.Credit("USDT", 100)
	_ = acc.Reserve("USDT", 60)
	err := acc.UseReserved("USDT", 70)
	assert.ErrorIs(suite.T(), err, account.ErrInsufficient)
}

func (suite *AccountTestSuite) TestReleaseReserved_Success() {
	acc, _ := account.NewAccount(suite.propsFaker, idObjValue.Uuid)
	_ = acc.Credit("BRL", 100)
	_ = acc.Reserve("BRL", 40)
	err := acc.ReleaseReserved("BRL", 30)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(90), acc.Balances["BRL"].Available)
	assert.Equal(suite.T(), int64(10), acc.Balances["BRL"].Reserved)
}

func (suite *AccountTestSuite) TestReleaseReserved_ErrorOnInsufficientReserved() {
	acc, _ := account.NewAccount(suite.propsFaker, idObjValue.Uuid)
	_ = acc.Credit("BRL", 100)
	_ = acc.Reserve("BRL", 40)
	err := acc.ReleaseReserved("BRL", 50)
	assert.ErrorIs(suite.T(), err, account.ErrInsufficient)
}

func (suite *AccountTestSuite) TestReserve_ErrorOnZeroOrNegativeAmount() {
	acc, _ := account.NewAccount(suite.propsFaker, idObjValue.Uuid)

	err := acc.Reserve("BTC", 0)
	suite.Equal(account.ErrInvalidParam, err)

	err = acc.Reserve("BTC", -10)
	suite.Equal(account.ErrInvalidParam, err)
}

func (suite *AccountTestSuite) TestUseReserved_ErrorOnNegativeAmount() {
	acc, _ := account.NewAccount(suite.propsFaker, idObjValue.Uuid)

	err := acc.UseReserved("BTC", -1)
	suite.Equal(account.ErrInvalidParam, err)
}

func (suite *AccountTestSuite) TestReleaseReserved_ErrorOnNegativeAmount() {
	acc, _ := account.NewAccount(suite.propsFaker, idObjValue.Uuid)

	err := acc.ReleaseReserved("BTC", -1)
	suite.Equal(account.ErrInvalidParam, err)
}

func (suite *AccountTestSuite) TestNewAccount_Valid() {
	props := account.AccountProps{Name: "user1"}
	acc, err := account.NewAccount(props, idObjValue.Uuid)
	suite.NoError(err)
	suite.NotNil(acc)
	suite.Equal("user1", acc.Name)
	suite.NotEmpty(acc.CreatedAt)
	suite.NotNil(acc.Balances)
}

func (suite *AccountTestSuite) TestNewAccount_InvalidName() {
	props := account.AccountProps{Name: ""}
	acc, err := account.NewAccount(props, idObjValue.Uuid)
	suite.ErrorIs(err, account.ErrInvalidParam)
	suite.Nil(acc)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}
