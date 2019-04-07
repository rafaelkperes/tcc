package gen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandUint64s(t *testing.T) {
	r := require.New(t)

	res := randUint64s(100, 5, 10)
	r.Len(res, 100)
	for _, val := range res {
		r.LessOrEqual(uint64(5), val)
		r.GreaterOrEqual(uint64(10), val)
	}
}
