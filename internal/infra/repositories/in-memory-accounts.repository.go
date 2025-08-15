package repositories

import (
	"sync"

	"github.com/juninhoitabh/clob-go/internal/domain/accounts"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type InMemoryAccountsRepository struct {
	mu       sync.Mutex
	accounts map[string]*accounts.Account
}

func NewInMemoryAccountsRepository() *InMemoryAccountsRepository {
	return &InMemoryAccountsRepository{
		accounts: make(map[string]*accounts.Account),
	}
}

func (i *InMemoryAccountsRepository) Create(id, name string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.accounts[id]; ok {
		return false
	}

	i.accounts[id] = &accounts.Account{
		ID:       id,
		Name:     name,
		Balances: make(map[string]*accounts.Balance),
	}

	return true
}

func (i *InMemoryAccountsRepository) Get(id string) (*accounts.Account, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	acct, ok := i.accounts[id]
	if !ok {
		return nil, shared.ErrNotFound
	}

	return acct, nil
}

func (i *InMemoryAccountsRepository) Credit(id, asset string, amount int64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	acct, ok := i.accounts[id]
	if !ok {
		return shared.ErrNotFound
	}

	return acct.Credit(asset, amount)
}

func (i *InMemoryAccountsRepository) Reserve(id, asset string, amount int64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	acct, ok := i.accounts[id]
	if !ok {
		return shared.ErrNotFound
	}

	return acct.Reserve(asset, amount)
}

func (i *InMemoryAccountsRepository) UseReserved(id, asset string, amount int64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	acct, ok := i.accounts[id]
	if !ok {
		return shared.ErrNotFound
	}

	return acct.UseReserved(asset, amount)
}

func (i *InMemoryAccountsRepository) ReleaseReserved(id, asset string, amount int64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	acct, ok := i.accounts[id]
	if !ok {
		return shared.ErrNotFound
	}

	return acct.ReleaseReserved(asset, amount)
}
