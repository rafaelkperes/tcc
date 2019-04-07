package gen

import (
	"math"
	"math/rand"
)

func bigInts(size int) []uint64 {
	return randUint64s(size, math.MaxUint8+1, math.MaxUint64)
}

func smallInts(size int) []uint64 {
	return randUint64s(size, 0, math.MaxUint8)
}

func allInts(size int) []uint64 {
	smallSize := int(float64(size) * (float64(math.MaxUint8) / float64(math.MaxUint64)))
	bigSize := size - smallSize

	s := smallInts(smallSize)
	b := bigInts(bigSize)
	res := append(b, s...)
	rand.Shuffle(len(res), func(i, j int) {
		t := res[i]
		res[i] = res[j]
		res[j] = t
	})
	return res
}

func randUint64s(total int, min uint64, max uint64) []uint64 {
	res := make([]uint64, total)
	for idx := range res {
		res[idx] = rand.Uint64()%(max-min) + min
	}
	return res
}
