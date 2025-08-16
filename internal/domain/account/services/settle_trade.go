package services

import (
	"fmt"

	"github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

func SettleTrade(
	repo account.AccountRepository,
	buyerID, sellerID string,
	base, quote string,
	price, qty int64,
) error {
	cost := shared.Mul(price, qty)

	if err := repo.UseReserved(buyerID, quote, cost); err != nil {
		return fmt.Errorf("buyer use reserved: %w", err)
	}
	if err := repo.Transfer(buyerID, base, qty); err != nil {
		return fmt.Errorf("transfer base to buyer: %w", err)
	}

	if err := repo.UseReserved(sellerID, base, qty); err != nil {
		return fmt.Errorf("seller use reserved: %w", err)
	}
	if err := repo.Transfer(sellerID, quote, cost); err != nil {
		return fmt.Errorf("transfer quote to seller: %w", err)
	}

	return nil
}
