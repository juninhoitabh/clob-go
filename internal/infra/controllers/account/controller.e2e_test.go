package account_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	httpAdapter "github.com/juninhoitabh/clob-go/internal/infra/http-client"
	httpServer "github.com/juninhoitabh/clob-go/internal/infra/http-server"
)

type (
	createInputDtoTest struct {
		AccountName string `json:"account_name"`
	}
	createOutputDtoTest struct {
		AccountId string `json:"account_id"`
		Status    string `json:"status"`
	}
	getAllByIdBalanceOutputDtoTest struct {
		Available int64 `json:"available"`
		Reserved  int64 `json:"reserved"`
	}
	getAllByIdOutputDtoTest struct {
		Balances  map[string]getAllByIdBalanceOutputDtoTest `json:"balances"`
		AccountID string                                    `json:"account_id"`
	}
	creditInputDtoTest struct {
		Asset  string `json:"asset"`
		Amount int64  `json:"amount"`
	}
	creditBalanceOutputDtoTest struct {
		Available int64 `json:"available"`
		Reserved  int64 `json:"reserved"`
	}
	creditOutputDtoTest struct {
		Balances  map[string]creditBalanceOutputDtoTest `json:"balances"`
		AccountID string                                `json:"account_id"`
	}
	AccountControllerTestSuite struct {
		suite.Suite
		httpClient    httpAdapter.HttpClient
		e2eTestHandle *httpServer.E2eTestHandle
		basePath      string
	}
)

func (suite *AccountControllerTestSuite) SetupTest() {
	suite.e2eTestHandle = httpServer.NewE2eTestHandle()
	suite.basePath = suite.e2eTestHandle.HttpServerTest.URL + "/api/v1/accounts"
	suite.httpClient = httpAdapter.NewDefaultHttpClient(10 * time.Second)
}

func (suite *AccountControllerTestSuite) TearDownSuite() {
	suite.e2eTestHandle.HttpServerTest.Close()
}

func (suite *AccountControllerTestSuite) TestCreate_Success() {
	t := suite.Suite.T()
	ctx := context.Background()

	inputBody := createInputDtoTest{AccountName: "e2e-alice"}

	headers := map[string]string{"Content-Type": "application/json"}
	res, err := suite.httpClient.Post(ctx, suite.basePath, inputBody, headers)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var out createOutputDtoTest
	err = json.Unmarshal(res.Body, &out)
	require.NoError(t, err)
	require.NotEmpty(t, out.AccountId)
	assert.Equal(t, "created", out.Status)
}

func (suite *AccountControllerTestSuite) TestCreate_MissingName_ReturnsBadRequest() {
	t := suite.Suite.T()

	inputBody := createInputDtoTest{AccountName: ""}
	dtoToSend, err := json.Marshal(inputBody)
	require.NoError(t, err)

	res, err := http.Post(suite.basePath, "application/json", bytes.NewReader(dtoToSend))
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func (suite *AccountControllerTestSuite) TestCreate_DuplicateName_ReturnsExists() {
	t := suite.Suite.T()

	in := createInputDtoTest{AccountName: "e2e-duplicate"}
	body, err := json.Marshal(in)
	require.NoError(t, err)

	res, err := http.Post(suite.basePath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)

	defer res.Body.Close()

	assert.True(t, res.StatusCode == http.StatusCreated || res.StatusCode == http.StatusOK)

	res2, err := http.Post(suite.basePath, "application/json", bytes.NewReader(body))
	require.NoError(t, err)

	defer res2.Body.Close()

	assert.Equal(t, http.StatusOK, res2.StatusCode)

	var out createOutputDtoTest
	err = json.NewDecoder(res2.Body).Decode(&out)
	require.NoError(t, err)
	assert.Equal(t, "exists", out.Status)
}

func (suite *AccountControllerTestSuite) TestCreate_InvalidJSON_ReturnsBadRequest() {
	t := suite.Suite.T()

	res, err := http.Post(suite.basePath, "application/json", bytes.NewBufferString("{invalid-json"))
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func (suite *AccountControllerTestSuite) TestGetAllById_Success() {
	t := suite.Suite.T()

	createInput := createInputDtoTest{AccountName: "get-account-test"}
	createBody, err := json.Marshal(createInput)
	require.NoError(t, err)

	createRes, err := http.Post(suite.basePath, "application/json", bytes.NewReader(createBody))
	require.NoError(t, err)
	defer createRes.Body.Close()

	assert.Equal(t, http.StatusCreated, createRes.StatusCode)

	var createOut createOutputDtoTest
	err = json.NewDecoder(createRes.Body).Decode(&createOut)
	require.NoError(t, err)

	creditInput := creditInputDtoTest{
		Asset:  "BTC",
		Amount: 1000,
	}
	creditBody, err := json.Marshal(creditInput)
	require.NoError(t, err)

	creditURL := suite.basePath + "/" + createOut.AccountId + "/credit"
	response, err := http.Post(creditURL, "application/json", bytes.NewReader(creditBody))
	require.NoError(t, err)

	defer response.Body.Close()

	getURL := suite.basePath + "/" + createOut.AccountId
	getRes, err := http.Get(getURL)
	require.NoError(t, err)

	defer getRes.Body.Close()

	assert.Equal(t, http.StatusOK, getRes.StatusCode)

	var getOut getAllByIdOutputDtoTest
	err = json.NewDecoder(getRes.Body).Decode(&getOut)
	require.NoError(t, err)

	assert.Equal(t, createOut.AccountId, getOut.AccountID)
	assert.NotNil(t, getOut.Balances)
	assert.Contains(t, getOut.Balances, "BTC")
	assert.Equal(t, int64(1000), getOut.Balances["BTC"].Available)
	assert.Equal(t, int64(0), getOut.Balances["BTC"].Reserved)
}

func (suite *AccountControllerTestSuite) TestGetAllById_NotFound() {
	t := suite.Suite.T()

	getURL := suite.basePath + "/non-existent-id"
	getRes, err := http.Get(getURL)
	require.NoError(t, err)

	defer getRes.Body.Close()

	assert.Equal(t, http.StatusNotFound, getRes.StatusCode)
}

func (suite *AccountControllerTestSuite) TestCredit_Success() {
	t := suite.Suite.T()

	createInput := createInputDtoTest{AccountName: "credit-account-test"}
	createBody, err := json.Marshal(createInput)
	require.NoError(t, err)

	createRes, err := http.Post(suite.basePath, "application/json", bytes.NewReader(createBody))
	require.NoError(t, err)
	defer createRes.Body.Close()

	assert.Equal(t, http.StatusCreated, createRes.StatusCode)

	var createOut createOutputDtoTest
	err = json.NewDecoder(createRes.Body).Decode(&createOut)
	require.NoError(t, err)

	creditInput := creditInputDtoTest{
		Asset:  "BTC",
		Amount: 100,
	}
	creditBody, err := json.Marshal(creditInput)
	require.NoError(t, err)

	creditURL := suite.basePath + "/" + createOut.AccountId + "/credit"
	creditRes, err := http.Post(creditURL, "application/json", bytes.NewReader(creditBody))
	require.NoError(t, err)

	defer creditRes.Body.Close()

	assert.Equal(t, http.StatusOK, creditRes.StatusCode)

	var creditOut creditOutputDtoTest
	err = json.NewDecoder(creditRes.Body).Decode(&creditOut)
	require.NoError(t, err)
	assert.Equal(t, createOut.AccountId, creditOut.AccountID)
	assert.NotNil(t, creditOut.Balances)
	assert.Contains(t, creditOut.Balances, "BTC")
	assert.Equal(t, int64(100), creditOut.Balances["BTC"].Available)
	assert.Equal(t, int64(0), creditOut.Balances["BTC"].Reserved)
}

func (suite *AccountControllerTestSuite) TestCredit_InvalidJSON() {
	t := suite.Suite.T()

	createInput := createInputDtoTest{AccountName: "credit-invalid-json"}
	createBody, err := json.Marshal(createInput)
	require.NoError(t, err)

	createRes, err := http.Post(suite.basePath, "application/json", bytes.NewReader(createBody))
	require.NoError(t, err)
	defer createRes.Body.Close()

	var createOut createOutputDtoTest
	err = json.NewDecoder(createRes.Body).Decode(&createOut)
	require.NoError(t, err)

	creditURL := suite.basePath + "/" + createOut.AccountId + "/credit"
	creditRes, err := http.Post(creditURL, "application/json", bytes.NewBufferString("{invalid-json"))
	require.NoError(t, err)

	defer creditRes.Body.Close()

	assert.Equal(t, http.StatusBadRequest, creditRes.StatusCode)
}

func (suite *AccountControllerTestSuite) TestCredit_InvalidAmount() {
	t := suite.Suite.T()

	createInput := createInputDtoTest{AccountName: "credit-invalid-amount"}
	createBody, err := json.Marshal(createInput)
	require.NoError(t, err)

	createRes, err := http.Post(suite.basePath, "application/json", bytes.NewReader(createBody))
	require.NoError(t, err)
	defer createRes.Body.Close()

	var createOut createOutputDtoTest
	err = json.NewDecoder(createRes.Body).Decode(&createOut)
	require.NoError(t, err)

	creditInput := creditInputDtoTest{
		Asset:  "BTC",
		Amount: -100,
	}
	creditBody, err := json.Marshal(creditInput)
	require.NoError(t, err)

	creditURL := suite.basePath + "/" + createOut.AccountId + "/credit"
	creditRes, err := http.Post(creditURL, "application/json", bytes.NewReader(creditBody))
	require.NoError(t, err)

	defer creditRes.Body.Close()

	assert.Equal(t, http.StatusBadRequest, creditRes.StatusCode)
}

func (suite *AccountControllerTestSuite) TestCredit_AccountNotFound() {
	t := suite.Suite.T()

	creditInput := creditInputDtoTest{
		Asset:  "BTC",
		Amount: 100,
	}
	creditBody, err := json.Marshal(creditInput)
	require.NoError(t, err)

	creditURL := suite.basePath + "/non-existent-id/credit"
	creditRes, err := http.Post(creditURL, "application/json", bytes.NewReader(creditBody))
	require.NoError(t, err)

	defer creditRes.Body.Close()

	assert.Equal(t, http.StatusNotFound, creditRes.StatusCode)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(AccountControllerTestSuite))
}
