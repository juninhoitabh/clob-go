package account

type (
	AccountSnapshot struct {
		Balances  map[string]Balance
		AccountID string
	}
)

type IAccountDAO interface {
	Snapshot(id string) (*AccountSnapshot, error)
}
