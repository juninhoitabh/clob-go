package account

type (
	AccountSnapshot struct {
		AccountID string
		Balances  map[string]Balance
	}
)

type IAccountDAO interface {
	Snapshot(id string) (*AccountSnapshot, error)
}
