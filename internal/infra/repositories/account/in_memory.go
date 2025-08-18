package repositories

import (
	"sync"

	domainAccount "github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type InMemoryAccountRepository struct {
	mu       sync.Mutex
	accounts map[string]*domainAccount.Account
}

func NewInMemoryAccountRepository() *InMemoryAccountRepository {
	return &InMemoryAccountRepository{
		accounts: make(map[string]*domainAccount.Account),
	}
}

func (i *InMemoryAccountRepository) Create(account *domainAccount.Account) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	id := account.GetID()

	if _, ok := i.accounts[id]; ok {
		return shared.ErrAlreadyExists
	}

	for _, acct := range i.accounts {
		if acct.Name == account.Name {
			return shared.ErrAlreadyExists
		}
	}

	i.accounts[id] = account

	return nil
}

func (i *InMemoryAccountRepository) Get(id string) (*domainAccount.Account, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	acct, ok := i.accounts[id]
	if !ok {
		return nil, shared.ErrNotFound
	}

	return acct, nil
}

func (i *InMemoryAccountRepository) Save(account *domainAccount.Account) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	id := account.GetID()
	i.accounts[id] = account

	return nil
}

func (i *InMemoryAccountRepository) AccountsMap() map[string]*domainAccount.Account {
	return i.accounts
}

func (i *InMemoryAccountRepository) Mutex() *sync.Mutex {
	return &i.mu
}
