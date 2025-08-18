package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/domain/book/services"
	"github.com/juninhoitabh/clob-go/internal/domain/order"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

type MatchOrderUnitTestSuite struct {
	suite.Suite
	book *book.Book
}

func (suite *MatchOrderUnitTestSuite) SetupTest() {
	props := book.BookProps{Instrument: "BTC/USDT"}
	b, _ := book.NewBook(props, idObjValue.Uuid)
	suite.book = b
}

func (suite *MatchOrderUnitTestSuite) TestMatchOrder_FullMatchBuy() {
	ask := &order.Order{
		AccountID: "seller1",
		Side:      order.Sell,
		Price:     100,
		Qty:       10,
		Remaining: 10,
	}
	suite.book.AddOrder(ask)

	buy := &order.Order{
		AccountID: "buyer1",
		Side:      order.Buy,
		Price:     100,
		Qty:       10,
		Remaining: 10,
	}

	report := services.MatchOrder(suite.book, buy)
	assert.Len(suite.T(), report.Trades, 1)
	trade := report.Trades[0]
	assert.Equal(suite.T(), int64(10), trade.Qty)
	assert.Equal(suite.T(), int64(100), trade.Price)
	assert.Equal(suite.T(), "buyer1", trade.Buyer)
	assert.Equal(suite.T(), "seller1", trade.Seller)
	assert.Equal(suite.T(), int64(0), buy.Remaining)
	assert.Equal(suite.T(), int64(0), ask.Remaining)
}

func (suite *MatchOrderUnitTestSuite) TestMatchOrder_PartialMatchBuy() {
	ask := &order.Order{
		AccountID: "seller2",
		Side:      order.Sell,
		Price:     100,
		Qty:       5,
		Remaining: 5,
	}
	suite.book.AddOrder(ask)

	buy := &order.Order{
		AccountID: "buyer2",
		Side:      order.Buy,
		Price:     100,
		Qty:       10,
		Remaining: 10,
	}

	report := services.MatchOrder(suite.book, buy)
	assert.Len(suite.T(), report.Trades, 1)
	trade := report.Trades[0]
	assert.Equal(suite.T(), int64(5), trade.Qty)
	assert.Equal(suite.T(), int64(100), trade.Price)
	assert.Equal(suite.T(), int64(5), buy.Remaining)
	assert.Equal(suite.T(), int64(0), ask.Remaining)
	// Buy order deve ser adicionada ao book
	assert.Contains(suite.T(), suite.book.Bids()[100].Orders, buy)
}

func (suite *MatchOrderUnitTestSuite) TestMatchOrder_NoMatchBuy() {
	buy := &order.Order{
		AccountID: "buyer3",
		Side:      order.Buy,
		Price:     99,
		Qty:       10,
		Remaining: 10,
	}
	report := services.MatchOrder(suite.book, buy)
	assert.Len(suite.T(), report.Trades, 0)
	assert.Equal(suite.T(), int64(10), buy.Remaining)
	assert.Contains(suite.T(), suite.book.Bids()[99].Orders, buy)
}

func (suite *MatchOrderUnitTestSuite) TestMatchOrder_FullMatchSell() {
	bid := &order.Order{
		AccountID: "buyer4",
		Side:      order.Buy,
		Price:     101,
		Qty:       8,
		Remaining: 8,
	}
	suite.book.AddOrder(bid)

	sell := &order.Order{
		AccountID: "seller4",
		Side:      order.Sell,
		Price:     100,
		Qty:       8,
		Remaining: 8,
	}

	report := services.MatchOrder(suite.book, sell)
	assert.Len(suite.T(), report.Trades, 1)
	trade := report.Trades[0]
	assert.Equal(suite.T(), int64(8), trade.Qty)
	assert.Equal(suite.T(), int64(101), trade.Price)
	assert.Equal(suite.T(), "buyer4", trade.Buyer)
	assert.Equal(suite.T(), "seller4", trade.Seller)
	assert.Equal(suite.T(), int64(0), sell.Remaining)
	assert.Equal(suite.T(), int64(0), bid.Remaining)
}

func (suite *MatchOrderUnitTestSuite) TestMatchOrder_MultiMatchBuy() {
	ask1 := &order.Order{
		AccountID: "seller5",
		Side:      order.Sell,
		Price:     100,
		Qty:       5,
		Remaining: 5,
	}
	ask2 := &order.Order{
		AccountID: "seller6",
		Side:      order.Sell,
		Price:     100,
		Qty:       7,
		Remaining: 7,
	}
	suite.book.AddOrder(ask1)
	suite.book.AddOrder(ask2)

	buy := &order.Order{
		AccountID: "buyer5",
		Side:      order.Buy,
		Price:     100,
		Qty:       10,
		Remaining: 10,
	}

	report := services.MatchOrder(suite.book, buy)
	assert.Len(suite.T(), report.Trades, 2)
	assert.Equal(suite.T(), int64(0), buy.Remaining)
	assert.Equal(suite.T(), int64(0), ask1.Remaining)
	assert.Equal(suite.T(), int64(2), ask2.Remaining)
}

func (suite *MatchOrderUnitTestSuite) TestMatchOrder_BuyOrderAddedWhenPartialFill() {
	b, _ := book.NewBook(book.BookProps{Instrument: "BTC/USDT"}, idObjValue.Uuid)
	ask := &order.Order{
		AccountID: "sellerX",
		Side:      order.Sell,
		Price:     100,
		Qty:       5,
		Remaining: 5,
	}
	b.AddOrder(ask)

	buy := &order.Order{
		AccountID: "buyerX",
		Side:      order.Buy,
		Price:     100,
		Qty:       10,
		Remaining: 10,
	}

	services.MatchOrder(b, buy)
	assert.Equal(suite.T(), int64(5), buy.Remaining)
	assert.Contains(suite.T(), b.Bids()[100].Orders, buy)
}

func (suite *MatchOrderUnitTestSuite) TestMatchOrder_SellOrderAddedWhenPartialFill() {
	b, _ := book.NewBook(book.BookProps{Instrument: "BTC/USDT"}, idObjValue.Uuid)
	bid := &order.Order{
		AccountID: "buyerY",
		Side:      order.Buy,
		Price:     100,
		Qty:       5,
		Remaining: 5,
	}
	b.AddOrder(bid)

	sell := &order.Order{
		AccountID: "sellerY",
		Side:      order.Sell,
		Price:     100,
		Qty:       10,
		Remaining: 10,
	}

	services.MatchOrder(b, sell)
	assert.Equal(suite.T(), int64(5), sell.Remaining)
	assert.Contains(suite.T(), b.Asks()[100].Orders, sell)
}

func (suite *MatchOrderUnitTestSuite) TestMatchOrder_SellPartialFillMakerRemains() {
	b, _ := book.NewBook(book.BookProps{Instrument: "BTC/USDT"}, idObjValue.Uuid)
	bid := &order.Order{
		AccountID: "buyerZ",
		Side:      order.Buy,
		Price:     100,
		Qty:       10,
		Remaining: 10,
	}
	b.AddOrder(bid)

	sell := &order.Order{
		AccountID: "sellerZ",
		Side:      order.Sell,
		Price:     100,
		Qty:       5,
		Remaining: 5,
	}

	services.MatchOrder(b, sell)
	assert.Equal(suite.T(), int64(5), bid.Remaining)
	assert.Equal(suite.T(), bid, b.Bids()[100].Orders[0])
}

func TestMatchOrderUnitTestSuite(t *testing.T) {
	suite.Run(t, new(MatchOrderUnitTestSuite))
}
