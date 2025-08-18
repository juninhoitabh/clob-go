package book_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/domain/order"
)

type PriceLevelUnitTestSuite struct {
	suite.Suite
}

func (suite *PriceLevelUnitTestSuite) TestNewPriceLevel() {
	pl := book.NewPriceLevel(100)
	assert.NotNil(suite.T(), pl)
	assert.Equal(suite.T(), int64(100), pl.Price)
	assert.Empty(suite.T(), pl.Orders)
}

func (suite *PriceLevelUnitTestSuite) TestTotalQty_Empty() {
	pl := book.NewPriceLevel(200)
	assert.Equal(suite.T(), int64(0), pl.TotalQty())
}

func (suite *PriceLevelUnitTestSuite) TestTotalQty_WithOrders() {
	pl := book.NewPriceLevel(300)
	pl.Orders = []*order.Order{
		{Remaining: 10},
		{Remaining: 20},
		{Remaining: 5},
	}
	assert.Equal(suite.T(), int64(35), pl.TotalQty())
}
