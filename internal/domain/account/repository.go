package account

type IAccountRepository interface {
	Create(account *Account) bool
	Save(account *Account) error
	Get(id string) (*Account, error)
}
