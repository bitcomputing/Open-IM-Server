package header

import (
	"context"
	"fmt"
)

type HeaderValuesType string

var HeaderValuesKey = HeaderValuesType("headerValues")

type HeaderValues struct {
	Token string
}

func GetHeaderValues(ctx context.Context) (*HeaderValues, error) {
	value, ok := ctx.Value(HeaderValuesKey).(*HeaderValues)
	if !ok {
		return nil, fmt.Errorf("type mismatch")
	}
	return value, nil
}
