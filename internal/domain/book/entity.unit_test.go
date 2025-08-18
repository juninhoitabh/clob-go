//go:build all || unit || domain

package book_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/domain/book/fakers"
	"github.com/juninhoitabh/clob-go/internal/domain/order"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

type BookUnitTestSuite struct {
	suite.Suite
	propsFaker book.BookProps
}

func (suite *BookUnitTestSuite) SetupTest() {
	suite.propsFaker = fakers.BookPropsFaker()
}

func (suite *BookUnitTestSuite) TestNewBook_Success() {
	b, err := book.NewBook(suite.propsFaker, idObjValue.Uuid)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), b)
	assert.Equal(suite.T(), suite.propsFaker.Instrument, b.Instrument)
	assert.NotNil(suite.T(), b.Bids())
	assert.NotNil(suite.T(), b.Asks())
}

func (suite *BookUnitTestSuite) TestNewBook_InvalidInstrument() {
	props := book.BookProps{Instrument: ""}
	b, err := book.NewBook(props, idObjValue.Uuid)
	assert.ErrorIs(suite.T(), err, book.ErrInvalidInstrumentBook)
	assert.Nil(suite.T(), b)
}

func (suite *BookUnitTestSuite) TestAddOrder_BidAndAsk() {
	b, _ := book.NewBook(suite.propsFaker, idObjValue.Uuid)
	bidOrder := &order.Order{Side: order.Buy, Price: 100}
	askOrder := &order.Order{Side: order.Sell, Price: 110}

	b.AddOrder(bidOrder)
	b.AddOrder(askOrder)

	assert.Contains(suite.T(), b.Bids(), int64(100))
	assert.Contains(suite.T(), b.Asks(), int64(110))
	assert.Equal(suite.T(), bidOrder, b.Bids()[100].Orders[0])
	assert.Equal(suite.T(), askOrder, b.Asks()[110].Orders[0])
}

func (suite *BookUnitTestSuite) TestRemoveOrder_RemovesOrderAndPriceLevel() {
	b, _ := book.NewBook(suite.propsFaker, idObjValue.Uuid)
	bidOrder := &order.Order{Side: order.Buy, Price: 100}
	b.AddOrder(bidOrder)

	assert.Contains(suite.T(), b.Bids(), int64(100))
	b.RemoveOrder(bidOrder)
	assert.NotContains(suite.T(), b.Bids(), int64(100))
}

func (suite *BookUnitTestSuite) TestBestBidAndBestAsk() {
	b, _ := book.NewBook(suite.propsFaker, idObjValue.Uuid)
	b.AddOrder(&order.Order{Side: order.Buy, Price: 100})
	b.AddOrder(&order.Order{Side: order.Buy, Price: 101})
	b.AddOrder(&order.Order{Side: order.Sell, Price: 110})
	b.AddOrder(&order.Order{Side: order.Sell, Price: 109})

	assert.Equal(suite.T(), int64(101), b.BestBid().Price)
	assert.Equal(suite.T(), int64(109), b.BestAsk().Price)
}

func (suite *BookUnitTestSuite) TestBidPricesAndAskPrices() {
	b, _ := book.NewBook(suite.propsFaker, idObjValue.Uuid)
	b.AddOrder(&order.Order{Side: order.Buy, Price: 100})
	b.AddOrder(&order.Order{Side: order.Buy, Price: 101})
	b.AddOrder(&order.Order{Side: order.Sell, Price: 110})
	b.AddOrder(&order.Order{Side: order.Sell, Price: 109})

	assert.ElementsMatch(suite.T(), []int64{101, 100}, b.BidPrices())
	assert.ElementsMatch(suite.T(), []int64{109, 110}, b.AskPrices())
}

func (suite *BookUnitTestSuite) TestBestBid_EmptyReturnsNil() {
	b, _ := book.NewBook(suite.propsFaker, idObjValue.Uuid)
	suite.Nil(b.BestBid())
}

func (suite *BookUnitTestSuite) TestBestAsk_EmptyReturnsNil() {
	b, _ := book.NewBook(suite.propsFaker, idObjValue.Uuid)
	suite.Nil(b.BestAsk())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(BookUnitTestSuite))
	suite.Run(t, new(PriceLevelUnitTestSuite))
}
