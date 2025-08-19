package order_test

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
	placeInputDtoTest struct {
		AccountID  string `json:"account_id"`
		Instrument string `json:"instrument"`
		Side       string `json:"side"`
		Price      int64  `json:"price"`
		Qty        int64  `json:"qty"`
	}
	placeTradeOutputDtoTest struct {
		TakerOrderID string `json:"taker_order_id"`
		MakerOrderID string `json:"maker_order_id"`
		BuyerID      string `json:"buyer_id"`
		SellerID     string `json:"seller_id"`
		Price        int64  `json:"price"`
		Qty          int64  `json:"qty"`
	}
	placeTradeReportOutputDtoTest struct {
		Trades []placeTradeOutputDtoTest `json:"trades"`
	}
	placeOutputDtoTest struct {
		Order  map[string]any                `json:"order"`
		Report placeTradeReportOutputDtoTest `json:"report"`
	}
	cancelOutputDtoTest struct {
		Order  map[string]any `json:"order"`
		Status string         `json:"status"`
	}
	OrderControllerTestSuite struct {
		suite.Suite
		e2eTestHandle *httpServer.E2eTestHandle
		basePath      string
		accountsPath  string
	}
)

func (suite *OrderControllerTestSuite) SetupTest() {
	suite.e2eTestHandle = httpServer.NewE2eTestHandle()
	suite.basePath = suite.e2eTestHandle.HttpServerTest.URL + "/api/v1/orders"
	suite.accountsPath = suite.e2eTestHandle.HttpServerTest.URL + "/api/v1/accounts"
}

func (suite *OrderControllerTestSuite) TearDownSuite() {
	suite.e2eTestHandle.HttpServerTest.Close()
}

func (suite *OrderControllerTestSuite) setupAccount(name string, asset string, amount int64) string {
	t := suite.Suite.T()

	createInput := map[string]string{"account_name": name}
	createBody, err := json.Marshal(createInput)
	require.NoError(t, err)

	createRes, err := http.Post(suite.accountsPath, "application/json", bytes.NewReader(createBody))
	require.NoError(t, err)
	defer createRes.Body.Close()

	var createOut map[string]string
	err = json.NewDecoder(createRes.Body).Decode(&createOut)
	require.NoError(t, err)

	accountID := createOut["account_id"]
	creditInput := map[string]interface{}{
		"asset":  asset,
		"amount": amount,
	}
	creditBody, err := json.Marshal(creditInput)
	require.NoError(t, err)

	creditURL := suite.accountsPath + "/" + accountID + "/credit"
	creditRes, err := http.Post(creditURL, "application/json", bytes.NewReader(creditBody))
	require.NoError(t, err)

	defer creditRes.Body.Close()

	return accountID
}

func (suite *OrderControllerTestSuite) TestPlace_Success() {
	t := suite.Suite.T()

	accountID := suite.setupAccount("order-test-account", "USDT", 100000)

	placeInput := placeInputDtoTest{
		AccountID:  accountID,
		Instrument: "BTC/USDT",
		Side:       "buy",
		Price:      50000,
		Qty:        1,
	}

	placeBody, err := json.Marshal(placeInput)
	require.NoError(t, err)

	placeRes, err := http.Post(suite.basePath, "application/json", bytes.NewReader(placeBody))
	require.NoError(t, err)
	defer placeRes.Body.Close()

	assert.Equal(t, http.StatusCreated, placeRes.StatusCode)

	var placeOut placeOutputDtoTest
	err = json.NewDecoder(placeRes.Body).Decode(&placeOut)
	require.NoError(t, err)

	assert.NotNil(t, placeOut.Order)
	assert.Equal(t, float64(50000), placeOut.Order["price"])
	assert.Equal(t, float64(1), placeOut.Order["qty"])
	assert.Equal(t, "buy", placeOut.Order["side"])
}

func (suite *OrderControllerTestSuite) TestPlace_InvalidJSON() {
	t := suite.Suite.T()

	res, err := http.Post(suite.basePath, "application/json", bytes.NewBufferString("{invalid-json"))
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func (suite *OrderControllerTestSuite) TestPlace_MissingFields() {
	t := suite.Suite.T()

	placeInput := placeInputDtoTest{
		Instrument: "BTC/USDT",
		Side:       "buy",
		Price:      50000,
		Qty:        1,
	}

	placeBody, err := json.Marshal(placeInput)
	require.NoError(t, err)

	placeRes, err := http.Post(suite.basePath, "application/json", bytes.NewReader(placeBody))
	require.NoError(t, err)
	defer placeRes.Body.Close()

	assert.Equal(t, http.StatusBadRequest, placeRes.StatusCode)
}

func (suite *OrderControllerTestSuite) TestPlace_InvalidParams() {
	t := suite.Suite.T()

	accountID := suite.setupAccount("order-invalid-params", "USDT", 100000)

	placeInput := placeInputDtoTest{
		AccountID:  accountID,
		Instrument: "BTC/USDT",
		Side:       "buy",
		Price:      -50000,
		Qty:        1,
	}

	placeBody, err := json.Marshal(placeInput)
	require.NoError(t, err)

	placeRes, err := http.Post(suite.basePath, "application/json", bytes.NewReader(placeBody))
	require.NoError(t, err)
	defer placeRes.Body.Close()

	assert.Equal(t, http.StatusBadRequest, placeRes.StatusCode)
}

func (suite *OrderControllerTestSuite) TestPlace_InsufficientBalance() {
	t := suite.Suite.T()

	accountID := suite.setupAccount("order-insufficient-balance", "USDT", 100)

	placeInput := placeInputDtoTest{
		AccountID:  accountID,
		Instrument: "BTC/USDT",
		Side:       "buy",
		Price:      5000,
		Qty:        1,
	}

	placeBody, err := json.Marshal(placeInput)
	require.NoError(t, err)

	placeRes, err := http.Post(suite.basePath, "application/json", bytes.NewReader(placeBody))
	require.NoError(t, err)
	defer placeRes.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, placeRes.StatusCode)

	bodyBytes, _ := io.ReadAll(placeRes.Body)
	assert.Contains(t, string(bodyBytes), "insufficient balance")
}

func (suite *OrderControllerTestSuite) TestPlace_AccountNotFound() {
	t := suite.Suite.T()

	placeInput := placeInputDtoTest{
		AccountID:  "non-existent-account",
		Instrument: "BTC/USDT",
		Side:       "buy",
		Price:      50000,
		Qty:        1,
	}

	placeBody, err := json.Marshal(placeInput)
	require.NoError(t, err)

	placeRes, err := http.Post(suite.basePath, "application/json", bytes.NewReader(placeBody))
	require.NoError(t, err)
	defer placeRes.Body.Close()

	assert.Equal(t, http.StatusNotFound, placeRes.StatusCode)
}

func (suite *OrderControllerTestSuite) TestCancel_Success() {
	t := suite.Suite.T()

	accountID := suite.setupAccount("cancel-test-account", "USDT", 100000)

	placeInput := placeInputDtoTest{
		AccountID:  accountID,
		Instrument: "BTC/USDT",
		Side:       "buy",
		Price:      50000,
		Qty:        1,
	}

	placeBody, err := json.Marshal(placeInput)
	require.NoError(t, err)

	placeRes, err := http.Post(suite.basePath, "application/json", bytes.NewReader(placeBody))
	require.NoError(t, err)
	defer placeRes.Body.Close()

	var placeOut placeOutputDtoTest
	err = json.NewDecoder(placeRes.Body).Decode(&placeOut)
	require.NoError(t, err)

	orderID := placeOut.Order["id"].(string)

	cancelURL := suite.basePath + "/" + orderID + "/cancel"
	cancelRes, err := http.Post(cancelURL, "application/json", nil)
	require.NoError(t, err)

	defer cancelRes.Body.Close()

	assert.Equal(t, http.StatusOK, cancelRes.StatusCode)

	var cancelOut cancelOutputDtoTest
	err = json.NewDecoder(cancelRes.Body).Decode(&cancelOut)
	require.NoError(t, err)

	assert.Equal(t, "canceled", cancelOut.Status)
	assert.Equal(t, orderID, cancelOut.Order["id"])
}

func (suite *OrderControllerTestSuite) TestCancel_OrderNotFound() {
	t := suite.Suite.T()

	cancelURL := suite.basePath + "/non-existent-order/cancel"
	cancelRes, err := http.Post(cancelURL, "application/json", nil)
	require.NoError(t, err)

	defer cancelRes.Body.Close()

	assert.Equal(t, http.StatusNotFound, cancelRes.StatusCode)
}

func (suite *OrderControllerTestSuite) TestCancel_MissingOrderID() {
	t := suite.Suite.T()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	cancelURL := suite.basePath + "//cancel"
	req, err := http.NewRequest("POST", cancelURL, nil)
	require.NoError(t, err)

	cancelRes, err := client.Do(req)
	require.NoError(t, err)
	defer cancelRes.Body.Close()

	assert.True(t, cancelRes.StatusCode == http.StatusMovedPermanently ||
		cancelRes.StatusCode == http.StatusBadRequest,
		"Status deve ser 301 ou 400, recebeu %d", cancelRes.StatusCode)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderControllerTestSuite))
}
