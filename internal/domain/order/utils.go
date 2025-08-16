package order

import (
	"errors"
	"strings"
)

var ErrInvalidSide = errors.New("invalid side")

func ParseSide(s string) (Side, error) {
	switch strings.ToLower(s) {
	case "buy":
		return Buy, nil
	case "sell":
		return Sell, nil
	default:
		return 0, ErrInvalidSide
	}
}
