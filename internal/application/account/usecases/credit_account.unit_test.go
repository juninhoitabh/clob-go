package usecases_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	accountUsecases "github.com/juninhoitabh/clob-go/internal/application/account/usecases"
	"github.com/juninhoitabh/clob-go/internal/application/account/usecases/fakers"
	domainAccount "github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/infra/repositories/account/mocks"
)

type CreditAccountUseCaseUnitTestSuite struct {
	suite.Suite
	inputFaker  accountUsecases.CreditAccountInput
	accountRepo *mocks.MockIAccountRepository
	ctrl        *gomock.Controller
	usecase     *accountUsecases.CreditAccountUseCase
}

func (suite *CreditAccountUseCaseUnitTestSuite) SetupTest() {
	suite.inputFaker = fakers.CreditAccountInputFaker()
	suite.ctrl = gomock.NewController(suite.T())
	suite.accountRepo = mocks.NewMockIAccountRepository(suite.ctrl)
	suite.usecase = accountUsecases.NewCreditAccountUseCase(suite.accountRepo)
}

func (suite *CreditAccountUseCaseUnitTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *CreditAccountUseCaseUnitTestSuite) TestExecute_Success() {
	input := suite.inputFaker
	account := &domainAccount.Account{}
	account.Balances = map[string]*domainAccount.Balance{
		input.Asset: {Available: 0, Reserved: 0},
	}

	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(gomock.Any()).Return(nil)

	err := suite.usecase.Execute(input)
	assert.NoError(suite.T(), err)
}

func (suite *CreditAccountUseCaseUnitTestSuite) TestExecute_GetError() {
	input := suite.inputFaker

	suite.accountRepo.EXPECT().Get(input.AccountID).Return(nil, errors.New("not found"))

	err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *CreditAccountUseCaseUnitTestSuite) TestExecute_SaveError() {
	input := suite.inputFaker
	account := &domainAccount.Account{}
	account.Balances = map[string]*domainAccount.Balance{
		input.Asset: {Available: 0, Reserved: 0},
	}

	suite.accountRepo.EXPECT().Get(input.AccountID).Return(account, nil)
	suite.accountRepo.EXPECT().Save(gomock.Any()).Return(errors.New("save error"))

	err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "save error")
}

func (suite *CreditAccountUseCaseUnitTestSuite) TestExecute_CreditError() {
	input := suite.inputFaker
	account := &domainAccount.Account{}
	account.Balances = map[string]*domainAccount.Balance{
		input.Asset: {Available: 0, Reserved: 0},
	}

	suite.accountRepo.EXPECT().Get(gomock.Any()).Return(account, nil)
	input.Amount = -10

	err := suite.usecase.Execute(input)
	assert.ErrorIs(suite.T(), err, domainAccount.ErrInvalidParam)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CreateAccountUseCaseUnitTestSuite))
	suite.Run(t, new(CreditAccountUseCaseUnitTestSuite))
}
