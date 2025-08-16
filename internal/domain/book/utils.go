package book

import (
	"fmt"
	"strings"
)

func SplitInstrument(inst string) (base string, quote string, err error) {
	parts := strings.Split(inst, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid instrument %q, expected BASE/QUOTE", inst)
	}
	return strings.ToUpper(parts[0]), strings.ToUpper(parts[1]), nil
}
