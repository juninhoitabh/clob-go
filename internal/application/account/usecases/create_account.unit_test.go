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

type CreateAccountUseCaseUnitTestSuite struct {
	suite.Suite
	inputFaker  accountUsecases.CreateAccountInput
	accountRepo *mocks.MockIAccountRepository
	ctrl        *gomock.Controller
	usecase     *accountUsecases.CreateAccountUseCase
}

func (suite *CreateAccountUseCaseUnitTestSuite) SetupTest() {
	suite.inputFaker = fakers.CreateAccountInputFaker()
	suite.ctrl = gomock.NewController(suite.T())
	suite.accountRepo = mocks.NewMockIAccountRepository(suite.ctrl)
	suite.usecase = accountUsecases.NewCreateAccountUseCase(suite.accountRepo)
}

func (suite *CreateAccountUseCaseUnitTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *CreateAccountUseCaseUnitTestSuite) TestExecute_Success() {
	input := suite.inputFaker

	// Cria entidade real para ser retornada
	account, _ := domainAccount.NewAccount(domainAccount.AccountProps{Name: input.AccountName}, "Uuid")

	suite.accountRepo.EXPECT().Create(gomock.Any()).Return(nil)

	output, err := suite.usecase.Execute(input)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), output)
	assert.Equal(suite.T(), account.Name, output.Name)
	assert.NotEmpty(suite.T(), output.ID)
}

func (suite *CreateAccountUseCaseUnitTestSuite) TestExecute_DomainError() {
	input := suite.inputFaker
	input.AccountName = "" // força erro de domínio

	output, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), output)
}

func (suite *CreateAccountUseCaseUnitTestSuite) TestExecute_RepoError() {
	input := suite.inputFaker

	suite.accountRepo.EXPECT().Create(gomock.Any()).Return(errors.New("repo error"))

	output, err := suite.usecase.Execute(input)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), output)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CreateAccountUseCaseUnitTestSuite))
}
