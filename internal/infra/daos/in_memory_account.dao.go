package daos

import (
	"sync"

	"github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type InMemoryAccountDAO struct {
	mu       *sync.Mutex
	accounts map[string]*account.Account
}

func NewInMemoryAccountDAO(mu *sync.Mutex, accounts map[string]*account.Account) *InMemoryAccountDAO {
	return &InMemoryAccountDAO{mu: mu, accounts: accounts}
}

func (dao *InMemoryAccountDAO) Snapshot(id string) (*account.AccountSnapshot, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	acct, ok := dao.accounts[id]
	if !ok {
		return nil, shared.ErrNotFound
	}

	out := &account.AccountSnapshot{
		AccountID: id,
		Balances:  make(map[string]account.Balance, len(acct.Balances)),
	}

	for asset, b := range acct.Balances {
		out.Balances[asset] = account.Balance{
			Available: b.Available,
			Reserved:  b.Reserved,
		}
	}

	return out, nil
}
