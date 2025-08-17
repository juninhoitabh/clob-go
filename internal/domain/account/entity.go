package account

import (
	"errors"
	"strings"
	"time"

	baseEntity "github.com/juninhoitabh/clob-go/internal/shared/domain/entities"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

var (
	ErrInsufficient = errors.New("insufficient balance")
	ErrInvalidParam = errors.New("invalid parameter")
)

type (
	AccountProps struct {
		Name string
	}
	Balance struct {
		Available int64
		Reserved  int64
	}
	Account struct {
		baseEntity.BaseEntity
		Name      string
		Balances  map[string]*Balance
		CreatedAt time.Time
	}
)

func (a *Account) Prepare(typeId idObjValue.TypeIdEnum) error {
	err := a.Validate()
	if err != nil {
		return err
	}

	a.NewBaseEntity("", typeId)

	a.CreatedAt = time.Now()

	a.Balances = make(map[string]*Balance)

	return nil
}

func (a *Account) Validate() error {
	if a.Name == "" {
		return ErrInvalidParam
	}

	return nil
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
	asset = strings.ToUpper(asset)

	bal, ok := a.Balances[asset]
	if !ok {
		bal = &Balance{}
		a.Balances[asset] = bal
	}
	return bal
}

func NewAccount(props AccountProps, typeId idObjValue.TypeIdEnum) (*Account, error) {
	account := Account{
		Name: props.Name,
	}

	err := account.Prepare(typeId)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
