package account

import "errors"

var (
	ErrInsufficient = errors.New("insufficient balance")
	ErrInvalidParam = errors.New("invalid parameter")
)

type Balance struct {
	Available int64
	Reserved  int64
}

type Account struct {
	ID       string // TODO: fazer objeto de valor
	Name     string // TODO: colocar como unica
	Balances map[string]*Balance
}

func (a *Account) Credit(asset string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidParam
	}
	bal := a.ensureBalance(asset)
	bal.Available += amount
	return nil
}

func (a *Account) Reserve(asset string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidParam
	}
	bal := a.ensureBalance(asset)
	if bal.Available < amount {
		return ErrInsufficient
	}
	bal.Available -= amount
	bal.Reserved += amount
	return nil
}

func (a *Account) UseReserved(asset string, amount int64) error {
	if amount < 0 {
		return ErrInvalidParam
	}
	bal := a.ensureBalance(asset)
	if bal.Reserved < amount {
		return ErrInsufficient
	}
	bal.Reserved -= amount
	return nil
}

func (a *Account) ReleaseReserved(asset string, amount int64) error {
	if amount < 0 {
		return ErrInvalidParam
	}
	bal := a.ensureBalance(asset)
	if bal.Reserved < amount {
		return ErrInsufficient
	}
	bal.Reserved -= amount
	bal.Available += amount
	return nil
}

func (a *Account) ensureBalance(asset string) *Balance {
	bal, ok := a.Balances[asset]
	if !ok {
		bal = &Balance{}
		a.Balances[asset] = bal
	}
	return bal
}
