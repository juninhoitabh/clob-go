package usecases

type (
	CreditAccountInput struct {
		AccountID string
		Asset     string
		Amount    int64
	}
	CreateAccountInput struct {
		AccountName string
	}
	CreateAccountOutput struct {
		ID   string
		Name string
	}
	ICreditAccountUseCase interface {
		Execute(input CreditAccountInput) error
	}
	ICreateAccountUseCase interface {
		Execute(input CreateAccountInput) (*CreateAccountOutput, error)
	}
)
