package services

import (
	"fmt"

	"github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

func SettleTrade(
	repo account.IAccountRepository,
	buyerID, sellerID string,
	base, quote string,
	price, qty int64,
) error {
	cost := shared.Mul(price, qty)

	buyerAcct, err := repo.Get(buyerID)
	if err != nil {
		return shared.ErrNotFound
	}

	sellerAcct, err := repo.Get(sellerID)
	if err != nil {
		return shared.ErrNotFound
	}

	if err := buyerAcct.UseReserved(quote, cost); err != nil {
		return fmt.Errorf("buyer use reserved: %w", err)
	}

	if err := buyerAcct.Credit(base, qty); err != nil {
		return fmt.Errorf("transfer base to buyer: %w", err)
	}

	if err := repo.Save(buyerAcct); err != nil {
		return err
	}

	if err := sellerAcct.UseReserved(base, qty); err != nil {
		return fmt.Errorf("seller use reserved: %w", err)
	}

	if err := sellerAcct.Credit(quote, cost); err != nil {
		return fmt.Errorf("transfer quote to seller: %w", err)
	}

	if err := repo.Save(sellerAcct); err != nil {
		return err
	}

	return nil
}
