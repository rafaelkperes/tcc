package data

import (
	"math"
)

const (
	defaultIntMin int64 = 0
	defaultIntMax int64 = math.MaxInt64
	defaultStrLen int64 = 100
)

type options struct {
	intMin int64
	intMax int64
	strLen int64
}

func newOptions() *options {
	return &options{
		intMin: defaultIntMin,
		intMax: defaultIntMax,
		strLen: defaultStrLen,
	}
}

// Option overrides behavior of Connect.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithInts(min, max int64) Option {
	return optionFunc(func(o *options) {
		o.intMin = min
		o.intMax = max
	})
}

func WithStringLength(l int64) Option {
	return optionFunc(func(o *options) {
		o.strLen = l
	})
}
