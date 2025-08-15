package services

import "github.com/juninhoitabh/clob-go/internal/domain/accounts"

func Transfer(from, to *accounts.Account, asset string, amount int64) error {
	if err := from.Reserve(asset, amount); err != nil {
		return err
	}

	if err := from.UseReserved(asset, amount); err != nil {
		return err
	}

	return to.Credit(asset, amount)
}
