//go:build all || unit || infra

package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	repositoriesBook "github.com/juninhoitabh/clob-go/internal/infra/repositories/book"
)

type InMemoryBookRepositoryE2eTestSuite struct {
	suite.Suite
	repo *repositoriesBook.InMemoryBookRepository
}

func (suite *InMemoryBookRepositoryE2eTestSuite) SetupTest() {
	suite.repo = repositoriesBook.NewInMemoryBookRepository()
}

func (suite *InMemoryBookRepositoryE2eTestSuite) TestSaveAndGetBook_Success() {
	book := &domainBook.Book{
		Instrument: "BTC/USDT",
	}
	err := suite.repo.SaveBook(book)
	suite.NoError(err)

	got, err := suite.repo.GetBook("BTC/USDT")
	suite.NoError(err)
	suite.Equal(book, got)
}

func (suite *InMemoryBookRepositoryE2eTestSuite) TestGetBook_NotFound() {
	got, err := suite.repo.GetBook("ETH/USDT")
	suite.NoError(err)
	suite.Nil(got)
}

func (suite *InMemoryBookRepositoryE2eTestSuite) TestSaveBook_Overwrite() {
	book1 := &domainBook.Book{Instrument: "BTC/USDT"}
	book2 := &domainBook.Book{Instrument: "BTC/USDT"}

	_ = suite.repo.SaveBook(book1)
	_ = suite.repo.SaveBook(book2)

	got, _ := suite.repo.GetBook("BTC/USDT")
	suite.Equal(book2, got)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryBookRepositoryE2eTestSuite))
}
