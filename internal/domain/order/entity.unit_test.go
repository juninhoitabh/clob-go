//go:build all || unit || domain

package order_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/juninhoitabh/clob-go/internal/domain/order"
	"github.com/juninhoitabh/clob-go/internal/domain/order/fakers"
	"github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

type OrderUnitTestSuite struct {
	suite.Suite
	propsFaker order.OrderProps
}

func (suite *OrderUnitTestSuite) SetupTest() {
	suite.propsFaker = fakers.OrderPropsFaker()
}

func (suite *OrderUnitTestSuite) TestNewOrder_Success() {
	props := order.OrderProps{
		AccountID:  "acc123",
		Instrument: "BTC/USDT",
		Side:       order.Buy,
		Price:      100,
		Qty:        10,
		Remaining:  10,
	}
	o, err := order.NewOrder(props, id.Uuid)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), o)
	assert.Equal(suite.T(), props.AccountID, o.AccountID)
	assert.Equal(suite.T(), props.Instrument, o.Instrument)
	assert.Equal(suite.T(), props.Side, o.Side)
	assert.Equal(suite.T(), props.Price, o.Price)
	assert.Equal(suite.T(), props.Qty, o.Qty)
	assert.Equal(suite.T(), props.Remaining, o.Remaining)
	assert.WithinDuration(suite.T(), time.Now(), o.CreatedAt, time.Second)
}

func (suite *OrderUnitTestSuite) TestNewOrder_InvalidAccountID() {
	props := suite.propsFaker
	props.AccountID = ""
	o, err := order.NewOrder(props, id.Uuid)
	assert.ErrorIs(suite.T(), err, order.ErrInvalidOrder)
	assert.Nil(suite.T(), o)
}

func (suite *OrderUnitTestSuite) TestNewOrder_InvalidInstrument() {
	props := suite.propsFaker
	props.Instrument = ""
	o, err := order.NewOrder(props, id.Uuid)
	assert.ErrorIs(suite.T(), err, order.ErrInvalidOrder)
	assert.Nil(suite.T(), o)
}

func (suite *OrderUnitTestSuite) TestNewOrder_InvalidSide() {
	props := suite.propsFaker
	props.Side = 0
	o, err := order.NewOrder(props, id.Uuid)
	assert.ErrorIs(suite.T(), err, order.ErrInvalidSideOrder)
	assert.Nil(suite.T(), o)
}

func (suite *OrderUnitTestSuite) TestNewOrder_InvalidPrice() {
	props := suite.propsFaker
	props.Price = 0
	o, err := order.NewOrder(props, id.Uuid)
	assert.ErrorIs(suite.T(), err, order.ErrInvalidOrder)
	assert.Nil(suite.T(), o)
}

func (suite *OrderUnitTestSuite) TestNewOrder_InvalidQty() {
	props := suite.propsFaker
	props.Qty = 0
	o, err := order.NewOrder(props, id.Uuid)
	assert.ErrorIs(suite.T(), err, order.ErrInvalidOrder)
	assert.Nil(suite.T(), o)
}

func (suite *OrderUnitTestSuite) TestNewOrder_InvalidRemainingNegative() {
	props := suite.propsFaker
	props.Remaining = -1
	o, err := order.NewOrder(props, id.Uuid)
	assert.ErrorIs(suite.T(), err, order.ErrInvalidOrder)
	assert.Nil(suite.T(), o)
}

func (suite *OrderUnitTestSuite) TestNewOrder_InvalidRemainingGreaterThanQty() {
	props := suite.propsFaker
	props.Remaining = props.Qty + 1
	o, err := order.NewOrder(props, id.Uuid)
	assert.ErrorIs(suite.T(), err, order.ErrInvalidOrder)
	assert.Nil(suite.T(), o)
}

func (suite *OrderUnitTestSuite) TestOrder_Public() {
	props := order.OrderProps{
		AccountID:  "acc123",
		Instrument: "BTC/USDT",
		Side:       order.Buy,
		Price:      100,
		Qty:        10,
		Remaining:  10,
	}
	o, err := order.NewOrder(props, id.Uuid)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), o)

	pub := o.Public()
	assert.Equal(suite.T(), o.ID, pub["id"])
	assert.Equal(suite.T(), o.AccountID, pub["account_id"])
	assert.Equal(suite.T(), o.Instrument, pub["instrument"])
	assert.Equal(suite.T(), o.Price, pub["price"])
	assert.Equal(suite.T(), o.Qty, pub["qty"])
	assert.Equal(suite.T(), o.Remaining, pub["remaining"])
	assert.Equal(suite.T(), o.CreatedAt.UTC().Format(time.RFC3339Nano), pub["created_at"])
	if o.Side == order.Buy {
		assert.Equal(suite.T(), "buy", pub["side"])
	} else {
		assert.Equal(suite.T(), "sell", pub["side"])
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderUnitTestSuite))
}
