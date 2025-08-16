package repositories

import (
	"sync"

	"github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type InMemoryAccountRepository struct {
	mu       sync.Mutex
	accounts map[string]*account.Account
}

func NewInMemoryAccountRepository() *InMemoryAccountRepository {
	return &InMemoryAccountRepository{
		accounts: make(map[string]*account.Account),
	}
}

func (i *InMemoryAccountRepository) Create(id, name string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.accounts[id]; ok {
		return false
	}

	i.accounts[id] = &account.Account{
		ID:       id,
		Name:     name,
		Balances: make(map[string]*account.Balance),
	}

	return true
}

func (i *InMemoryAccountRepository) Get(id string) (*account.Account, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	acct, ok := i.accounts[id]
	if !ok {
		return nil, shared.ErrNotFound
	}

	return acct, nil
}

func (i *InMemoryAccountRepository) Credit(id, asset string, amount int64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	acct, ok := i.accounts[id]
	if !ok {
		return shared.ErrNotFound
	}

	return acct.Credit(asset, amount)
}

func (i *InMemoryAccountRepository) Reserve(id, asset string, amount int64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	acct, ok := i.accounts[id]
	if !ok {
		return shared.ErrNotFound
	}

	return acct.Reserve(asset, amount)
}

func (i *InMemoryAccountRepository) UseReserved(id, asset string, amount int64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	acct, ok := i.accounts[id]
	if !ok {
		return shared.ErrNotFound
	}

	return acct.UseReserved(asset, amount)
}

func (i *InMemoryAccountRepository) ReleaseReserved(id, asset string, amount int64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	acct, ok := i.accounts[id]
	if !ok {
		return shared.ErrNotFound
	}

	return acct.ReleaseReserved(asset, amount)
}

func (i *InMemoryAccountRepository) AccountsMap() map[string]*account.Account {
	return i.accounts
}

func (i *InMemoryAccountRepository) Mutex() *sync.Mutex {
	return &i.mu
}
