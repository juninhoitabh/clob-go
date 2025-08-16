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
