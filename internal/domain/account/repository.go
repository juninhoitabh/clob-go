package account

type AccountsRepository interface {
	Create(id, name string) bool
	Get(id string) (*Account, error)
	Credit(id, asset string, amount int64) error
	Reserve(id, asset string, amount int64) error
	UseReserved(id, asset string, amount int64) error
	ReleaseReserved(id, asset string, amount int64) error
}
