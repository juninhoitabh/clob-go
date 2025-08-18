//go:build all || unit || domain

package book_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/juninhoitabh/clob-go/internal/domain/book"
)

func TestSplitInstrument_Valid(t *testing.T) {
	base, quote, err := book.SplitInstrument("btc/usdt")
	assert.NoError(t, err)
	assert.Equal(t, "BTC", base)
	assert.Equal(t, "USDT", quote)
}

func TestSplitInstrument_InvalidFormat(t *testing.T) {
	base, quote, err := book.SplitInstrument("btcusdt")
	assert.Error(t, err)
	assert.Empty(t, base)
	assert.Empty(t, quote)
	assert.Contains(t, err.Error(), "invalid instrument")
}

func TestSplitInstrument_ExtraParts(t *testing.T) {
	base, quote, err := book.SplitInstrument("btc/usdt/extra")
	assert.Error(t, err)
	assert.Empty(t, base)
	assert.Empty(t, quote)
	assert.Contains(t, err.Error(), "invalid instrument")
}

func TestSplitInstrument_UpperCase(t *testing.T) {
	base, quote, err := book.SplitInstrument("Eth/Brl")
	assert.NoError(t, err)
	assert.Equal(t, "ETH", base)
	assert.Equal(t, "BRL", quote)
}
