//go:build all || e2e || infra

package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	domainOrder "github.com/juninhoitabh/clob-go/internal/domain/order"
	repositoriesOrder "github.com/juninhoitabh/clob-go/internal/infra/repositories/order"
)

type InMemoryOrderRepositoryE2ETestSuite struct {
	suite.Suite
	repo *repositoriesOrder.InMemoryOrderRepository
}

func (suite *InMemoryOrderRepositoryE2ETestSuite) SetupTest() {
	suite.repo = repositoriesOrder.NewInMemoryOrderRepository()
}

func (suite *InMemoryOrderRepositoryE2ETestSuite) TestSaveAndGetOrder_Success() {
	order := &domainOrder.Order{}
	order.ID.ID = "order1"
	err := suite.repo.SaveOrder(order)
	assert.NoError(suite.T(), err)

	got, err := suite.repo.GetOrder("order1")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), order, got)
}

func (suite *InMemoryOrderRepositoryE2ETestSuite) TestGetOrder_NotFound() {
	got, err := suite.repo.GetOrder("unknown")
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), got)
}

func (suite *InMemoryOrderRepositoryE2ETestSuite) TestRemoveOrder() {
	order := &domainOrder.Order{}
	order.ID.ID = "order2"
	_ = suite.repo.SaveOrder(order)

	err := suite.repo.RemoveOrder("order2")
	assert.NoError(suite.T(), err)

	got, _ := suite.repo.GetOrder("order2")
	assert.Nil(suite.T(), got)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryOrderRepositoryE2ETestSuite))
}
