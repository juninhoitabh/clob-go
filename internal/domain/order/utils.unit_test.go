package order_test

import (
	"testing"

	"github.com/juninhoitabh/clob-go/internal/domain/order"
	"github.com/stretchr/testify/assert"
)

func TestParseSide_Buy(t *testing.T) {
	side, err := order.ParseSide("buy")
	assert.NoError(t, err)
	assert.Equal(t, order.Buy, side)
}

func TestParseSide_Sell(t *testing.T) {
	side, err := order.ParseSide("sell")
	assert.NoError(t, err)
	assert.Equal(t, order.Sell, side)
}

func TestParseSide_CaseInsensitive(t *testing.T) {
	side, err := order.ParseSide("BUY")
	assert.NoError(t, err)
	assert.Equal(t, order.Buy, side)

	side, err = order.ParseSide("SeLl")
	assert.NoError(t, err)
	assert.Equal(t, order.Sell, side)
}

func TestParseSide_Invalid(t *testing.T) {
	side, err := order.ParseSide("hold")
	assert.ErrorIs(t, err, order.ErrInvalidSide)
	assert.Equal(t, order.Side(0), side)
}
