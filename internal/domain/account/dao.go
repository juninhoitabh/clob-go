package account

type (
	AccountSnapshot struct {
		AccountID string
		Balances  map[string]Balance
	}
)

type AccountDAO interface {
	Snapshot(id string) (*AccountSnapshot, error)
}
