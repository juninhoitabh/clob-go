//go:build all || e2e || infra

package daos_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/juninhoitabh/clob-go/internal/domain/account"
	daosAccount "github.com/juninhoitabh/clob-go/internal/infra/daos/account"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type InMemoryAccountDAOE2ETestSuite struct {
	suite.Suite
	dao *daosAccount.InMemoryAccountDAO
}

func (suite *InMemoryAccountDAOE2ETestSuite) SetupTest() {
	mu := &sync.Mutex{}
	accounts := map[string]*account.Account{
		"acc1": {
			Balances: map[string]*account.Balance{
				"BTC":  {Available: 10, Reserved: 2},
				"USDT": {Available: 1000, Reserved: 0},
			},
		},
	}
	suite.dao = daosAccount.NewInMemoryAccountDAO(mu, accounts)
}

func (suite *InMemoryAccountDAOE2ETestSuite) TestSnapshot_Success() {
	snap, err := suite.dao.Snapshot("acc1")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), snap)
	assert.Equal(suite.T(), "acc1", snap.AccountID)
	assert.Equal(suite.T(), int64(10), snap.Balances["BTC"].Available)
	assert.Equal(suite.T(), int64(2), snap.Balances["BTC"].Reserved)
	assert.Equal(suite.T(), int64(1000), snap.Balances["USDT"].Available)
	assert.Equal(suite.T(), int64(0), snap.Balances["USDT"].Reserved)
}

func (suite *InMemoryAccountDAOE2ETestSuite) TestSnapshot_NotFound() {
	snap, err := suite.dao.Snapshot("unknown")
	assert.ErrorIs(suite.T(), err, shared.ErrNotFound)
	assert.Nil(suite.T(), snap)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryAccountDAOE2ETestSuite))
}
