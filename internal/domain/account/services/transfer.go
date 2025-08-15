package services

import "github.com/juninhoitabh/clob-go/internal/domain/account"

func Transfer(from, to *account.Account, asset string, amount int64) error {
	if err := from.Reserve(asset, amount); err != nil {
		return err
	}

	if err := from.UseReserved(asset, amount); err != nil {
		return err
	}

	return to.Credit(asset, amount)
}

// TODO: Era assim, ver como vai ficar
// func (i *InMemoryAccountsRepository) Transfer(receiver, asset string, amount int64) error {
// 	i.mu.Lock()
// 	defer i.mu.Unlock()

// 	acct, ok := i.accounts[receiver]
// 	if !ok {
// 		return shared.ErrNotFound
// 	}

// 	return acct.Credit(asset, amount)
// }
