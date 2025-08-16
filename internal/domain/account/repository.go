package account

type IAccountRepository interface {
	Create(account *Account) error
	Save(account *Account) error
	Get(id string) (*Account, error)
}
