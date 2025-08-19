package book_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	httpServer "github.com/juninhoitabh/clob-go/internal/infra/http-server"
)

type (
	getByInstrumentLevelOutputDtoTest struct {
		Price int64 `json:"price"`
		Qty   int64 `json:"qty"`
	}
	getByInstrumentOutputDtoTest struct {
		Instrument string                              `json:"instrument"`
		Bids       []getByInstrumentLevelOutputDtoTest `json:"bids"`
		Asks       []getByInstrumentLevelOutputDtoTest `json:"asks"`
	}
	BookControllerTestSuite struct {
		suite.Suite
		e2eTestHandle *httpServer.E2eTestHandle
		basePath      string
	}
)

func (suite *BookControllerTestSuite) SetupTest() {
	suite.e2eTestHandle = httpServer.NewE2eTestHandle()
	suite.basePath = suite.e2eTestHandle.HttpServerTest.URL + "/api/v1/books"
}

func (suite *BookControllerTestSuite) TearDownSuite() {
	suite.e2eTestHandle.HttpServerTest.Close()
}

func (suite *BookControllerTestSuite) TestGet_Success() {
	t := suite.Suite.T()

	accountsPath := suite.e2eTestHandle.HttpServerTest.URL + "/api/v1/accounts"
	createAccountInput := map[string]string{"account_name": "book-test-account"}
	createAccountBody, err := json.Marshal(createAccountInput)
	require.NoError(t, err)

	createAccountRes, err := http.Post(accountsPath, "application/json", bytes.NewReader(createAccountBody))
	require.NoError(t, err)
	defer createAccountRes.Body.Close()

	var createAccountOut map[string]string
	err = json.NewDecoder(createAccountRes.Body).Decode(&createAccountOut)
	require.NoError(t, err)

	accountID := createAccountOut["account_id"]
	creditInput := map[string]interface{}{
		"asset":  "USDT",
		"amount": 60000,
	}
	creditBody, err := json.Marshal(creditInput)
	require.NoError(t, err)

	creditURL := accountsPath + "/" + accountID + "/credit"
	creditRes, err := http.Post(creditURL, "application/json", bytes.NewReader(creditBody))
	require.NoError(t, err)

	defer creditRes.Body.Close()

	assert.Equal(t, http.StatusOK, creditRes.StatusCode)

	ordersPath := suite.e2eTestHandle.HttpServerTest.URL + "/api/v1/orders"
	createOrderInput := map[string]interface{}{
		"account_id": accountID,
		"instrument": "BTC/USDT",
		"side":       "buy",
		"qty":        1,
		"price":      50000,
	}
	createOrderBody, err := json.Marshal(createOrderInput)
	require.NoError(t, err)

	orderRes, err := http.Post(ordersPath, "application/json", bytes.NewReader(createOrderBody))
	require.NoError(t, err)
	defer orderRes.Body.Close()

	if orderRes.StatusCode != http.StatusCreated && orderRes.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(orderRes.Body)
		t.Logf("Erro ao criar ordem: %d - %s", orderRes.StatusCode, string(bodyBytes))
		t.FailNow()
	}

	getBookURL := suite.basePath + "?instrument=BTC/USDT"
	getBookRes, err := http.Get(getBookURL)
	require.NoError(t, err)

	defer getBookRes.Body.Close()

	assert.Equal(t, http.StatusOK, getBookRes.StatusCode)

	if getBookRes.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(getBookRes.Body)
		t.Logf("Erro ao buscar book: %d - %s", getBookRes.StatusCode, string(bodyBytes))
	}

	var bookOut getByInstrumentOutputDtoTest
	err = json.NewDecoder(getBookRes.Body).Decode(&bookOut)
	require.NoError(t, err)

	assert.Equal(t, "BTC/USDT", bookOut.Instrument)
	assert.NotNil(t, bookOut.Bids)
	assert.NotNil(t, bookOut.Asks)

	if assert.Len(t, bookOut.Bids, 1) {
		assert.Equal(t, int64(50000), bookOut.Bids[0].Price)
		assert.Equal(t, int64(1), bookOut.Bids[0].Qty)
	}
}

func (suite *BookControllerTestSuite) TestGet_EmptyInstrument() {
	t := suite.Suite.T()

	getBookURL := suite.basePath
	getBookRes, err := http.Get(getBookURL)
	require.NoError(t, err)

	defer getBookRes.Body.Close()

	assert.Equal(t, http.StatusBadRequest, getBookRes.StatusCode)
}

func (suite *BookControllerTestSuite) TestGet_InstrumentNotFound() {
	t := suite.Suite.T()

	getBookURL := suite.basePath + "?instrument=ETH/DOGE"
	getBookRes, err := http.Get(getBookURL)
	require.NoError(t, err)

	defer getBookRes.Body.Close()

	assert.Equal(t, http.StatusNotFound, getBookRes.StatusCode)
}

func (suite *BookControllerTestSuite) TestGet_CaseInsensitiveInstrument() {
	t := suite.Suite.T()

	accountsPath := suite.e2eTestHandle.HttpServerTest.URL + "/api/v1/accounts"
	createAccountInput := map[string]string{"account_name": "case-insensitive-test"}
	createAccountBody, err := json.Marshal(createAccountInput)
	require.NoError(t, err)

	createAccountRes, err := http.Post(accountsPath, "application/json", bytes.NewReader(createAccountBody))
	require.NoError(t, err)
	defer createAccountRes.Body.Close()

	var createAccountOut map[string]string
	err = json.NewDecoder(createAccountRes.Body).Decode(&createAccountOut)
	require.NoError(t, err)

	accountID := createAccountOut["account_id"]
	creditInput := map[string]interface{}{
		"asset":  "BRL",
		"amount": 60000,
	}
	creditBody, err := json.Marshal(creditInput)
	require.NoError(t, err)

	creditURL := accountsPath + "/" + accountID + "/credit"
	response, err := http.Post(creditURL, "application/json", bytes.NewReader(creditBody))
	require.NoError(t, err)

	defer response.Body.Close()

	ordersPath := suite.e2eTestHandle.HttpServerTest.URL + "/api/v1/orders"
	createOrderInput := map[string]interface{}{
		"account_id": accountID,
		"instrument": "BTC/BRL",
		"side":       "buy",
		"qty":        1,
		"price":      50000,
	}
	createOrderBody, err := json.Marshal(createOrderInput)
	require.NoError(t, err)

	response, err = http.Post(ordersPath, "application/json", bytes.NewReader(createOrderBody))
	require.NoError(t, err)

	defer response.Body.Close()

	getBookURL := suite.basePath + "?instrument=btc/brl"
	getBookRes, err := http.Get(getBookURL)
	require.NoError(t, err)

	defer getBookRes.Body.Close()

	assert.Equal(t, http.StatusOK, getBookRes.StatusCode)

	if getBookRes.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(getBookRes.Body)
		t.Logf("Erro ao buscar book: %d - %s", getBookRes.StatusCode, string(bodyBytes))
		t.FailNow()
	}

	var bookOut getByInstrumentOutputDtoTest
	err = json.NewDecoder(getBookRes.Body).Decode(&bookOut)
	require.NoError(t, err)

	assert.Equal(t, "BTC/BRL", bookOut.Instrument)
}

func (suite *BookControllerTestSuite) TestGet_WithAsks() {
	t := suite.Suite.T()

	accountsPath := suite.e2eTestHandle.HttpServerTest.URL + "/api/v1/accounts"
	createAccountInput := map[string]string{"account_name": "ask-test-account"}
	createAccountBody, err := json.Marshal(createAccountInput)
	require.NoError(t, err)

	createAccountRes, err := http.Post(accountsPath, "application/json", bytes.NewReader(createAccountBody))
	require.NoError(t, err)
	defer createAccountRes.Body.Close()

	var createAccountOut map[string]string
	err = json.NewDecoder(createAccountRes.Body).Decode(&createAccountOut)
	require.NoError(t, err)

	accountID := createAccountOut["account_id"]
	creditInput := map[string]interface{}{
		"asset":  "BTC",
		"amount": 5,
	}
	creditBody, err := json.Marshal(creditInput)
	require.NoError(t, err)

	creditURL := accountsPath + "/" + accountID + "/credit"
	creditRes, err := http.Post(creditURL, "application/json", bytes.NewReader(creditBody))
	require.NoError(t, err)

	defer creditRes.Body.Close()

	assert.Equal(t, http.StatusOK, creditRes.StatusCode)

	ordersPath := suite.e2eTestHandle.HttpServerTest.URL + "/api/v1/orders"
	createOrderInput := map[string]interface{}{
		"account_id": accountID,
		"instrument": "BTC/USDT",
		"side":       "sell",
		"qty":        2,
		"price":      55000,
	}
	createOrderBody, err := json.Marshal(createOrderInput)
	require.NoError(t, err)

	orderRes, err := http.Post(ordersPath, "application/json", bytes.NewReader(createOrderBody))
	require.NoError(t, err)
	defer orderRes.Body.Close()

	if orderRes.StatusCode != http.StatusCreated && orderRes.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(orderRes.Body)
		t.Logf("Erro ao criar ordem de venda: %d - %s", orderRes.StatusCode, string(bodyBytes))
		t.FailNow()
	}

	getBookURL := suite.basePath + "?instrument=BTC/USDT"
	getBookRes, err := http.Get(getBookURL)
	require.NoError(t, err)

	defer getBookRes.Body.Close()

	assert.Equal(t, http.StatusOK, getBookRes.StatusCode)

	if getBookRes.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(getBookRes.Body)
		t.Logf("Erro ao buscar book: %d - %s", getBookRes.StatusCode, string(bodyBytes))
		t.FailNow()
	}

	var bookOut getByInstrumentOutputDtoTest
	err = json.NewDecoder(getBookRes.Body).Decode(&bookOut)
	require.NoError(t, err)

	assert.Equal(t, "BTC/USDT", bookOut.Instrument)
	assert.NotNil(t, bookOut.Asks)

	if assert.Len(t, bookOut.Asks, 1) {
		assert.Equal(t, int64(55000), bookOut.Asks[0].Price)
		assert.Equal(t, int64(2), bookOut.Asks[0].Qty)
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(BookControllerTestSuite))
}
