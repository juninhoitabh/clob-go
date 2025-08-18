//go:build all || e2e || infra

package repositories_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	domainAccount "github.com/juninhoitabh/clob-go/internal/domain/account"
	repositoriesAccount "github.com/juninhoitabh/clob-go/internal/infra/repositories/account"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type InMemoryAccountRepositoryE2ETestSuite struct {
	suite.Suite
	repo *repositoriesAccount.InMemoryAccountRepository
}

func (suite *InMemoryAccountRepositoryE2ETestSuite) SetupTest() {
	suite.repo = repositoriesAccount.NewInMemoryAccountRepository()
}

func (suite *InMemoryAccountRepositoryE2ETestSuite) TestCreateAndGet_Success() {
	account, _ := domainAccount.NewAccount(domainAccount.AccountProps{Name: "Alice"}, "acc1")
	err := suite.repo.Create(account)
	assert.NoError(suite.T(), err)

	got, err := suite.repo.Get(account.GetID())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), account, got)
}

func (suite *InMemoryAccountRepositoryE2ETestSuite) TestCreate_AlreadyExistsByID() {
	account, _ := domainAccount.NewAccount(domainAccount.AccountProps{Name: "Bob"}, "acc2")
	_ = suite.repo.Create(account)
	err := suite.repo.Create(account)
	assert.ErrorIs(suite.T(), err, shared.ErrAlreadyExists)
}

func (suite *InMemoryAccountRepositoryE2ETestSuite) TestCreate_AlreadyExistsByName() {
	account1, _ := domainAccount.NewAccount(domainAccount.AccountProps{Name: "Carol"}, "acc3")
	account2, _ := domainAccount.NewAccount(domainAccount.AccountProps{Name: "Carol"}, "acc4")
	_ = suite.repo.Create(account1)
	err := suite.repo.Create(account2)
	assert.ErrorIs(suite.T(), err, shared.ErrAlreadyExists)
}

func (suite *InMemoryAccountRepositoryE2ETestSuite) TestGet_NotFound() {
	got, err := suite.repo.Get("unknown")
	assert.ErrorIs(suite.T(), err, shared.ErrNotFound)
	assert.Nil(suite.T(), got)
}

func (suite *InMemoryAccountRepositoryE2ETestSuite) TestSave_UpdatesAccount() {
	account, _ := domainAccount.NewAccount(domainAccount.AccountProps{Name: "Dave"}, "Uuid")
	_ = suite.repo.Create(account)
	account.Name = "DaveUpdated"
	err := suite.repo.Save(account)
	assert.NoError(suite.T(), err)

	got, _ := suite.repo.Get(account.GetID())
	assert.Equal(suite.T(), "DaveUpdated", got.Name)
}

func (suite *InMemoryAccountRepositoryE2ETestSuite) TestCreate_DuplicateName() {
	account1, _ := domainAccount.NewAccount(domainAccount.AccountProps{Name: "Alice"}, "Uuid")
	account2, _ := domainAccount.NewAccount(domainAccount.AccountProps{Name: "Alice"}, "Uuid")
	err := suite.repo.Create(account1)
	assert.NoError(suite.T(), err)
	err = suite.repo.Create(account2)
	assert.ErrorIs(suite.T(), err, shared.ErrAlreadyExists)
}

func (suite *InMemoryAccountRepositoryE2ETestSuite) TestAccountsMap() {
	account, _ := domainAccount.NewAccount(domainAccount.AccountProps{Name: "Bob"}, "Uuid")
	_ = suite.repo.Create(account)
	accountsMap := suite.repo.AccountsMap()
	assert.Contains(suite.T(), accountsMap, account.GetID())
	assert.Equal(suite.T(), account, accountsMap[account.GetID()])
}

func (suite *InMemoryAccountRepositoryE2ETestSuite) TestMutex() {
	mutex := suite.repo.Mutex()
	assert.IsType(suite.T(), &sync.Mutex{}, mutex)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryAccountRepositoryE2ETestSuite))
}
