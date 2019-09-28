package data

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"

	"github.com/sirupsen/logrus"
)

// Using fixed seed to always reproduce the same scenario.
// var src = rand.NewSource(time.Now().UnixNano())
var src = rand.NewSource(88)

func Create(typ Type, total, min, max, strlen int64) (Data, error) {
	switch typ {
	case TypeInt:
		return CreateInts(total, min, max), nil
	case TypeFloat:
		return CreateFloats(total), nil
	case TypeString:
		return CreateStrings(total, strlen), nil
	case TypeObject:
		return CreateObjects(total), nil
	default:
		return nil, fmt.Errorf("unknown type: %v", typ)
	}
}

func Unmarshal(d []byte, typ Type, format Format) (Data, error) {
	switch format {
	case FormatJSON:
		switch typ {
		case TypeInt:
			var v Ints
			err := json.Unmarshal(d, &v)
			return v, err
		case TypeFloat:
			var v Floats
			err := json.Unmarshal(d, &v)
			return v, err
		case TypeString:
			var v Strings
			err := json.Unmarshal(d, &v)
			return v, err
		case TypeObject:
			var v Objects
			err := json.Unmarshal(d, &v)
			return v, err
		default:
			return nil, fmt.Errorf("unknown type: %v", typ)
		}
	default:
		return nil, fmt.Errorf("unknown format: %v", format)
	}
}

func CreateInts(total, min, max int64) Ints {
	r := make(Ints, total)
	for idx := range r {
		r[idx] = genInt(min, max)
	}
	logrus.WithFields(map[string]interface{}{
		"typ":    TypeInt,
		"size":   total,
		"intmin": min,
		"intmax": max,
		"head":   r[:5],
	}).Debug("created ints")
	return r
}

func genInt(min, max int64) int64 {
	v := src.Int63()%(max-min) + min
	if src.Int63()%2 == 0 {
		return -1 * v
	}
	return v
}

func CreateFloats(total int64) Floats {
	rnd := rand.New(src)
	r := make(Floats, total)
	for idx := range r {
		r[idx] = rnd.NormFloat64()
	}
	logrus.WithFields(map[string]interface{}{
		"typ":  TypeInt,
		"size": total,
		"head": r[:5],
	}).Debug("created floats")
	return r
}

func CreateStrings(total, length int64) Strings {
	r := make(Strings, total)
	for idx := range r {
		r[idx] = genString(length)
	}

	logrus.WithFields(map[string]interface{}{
		"typ":    TypeString,
		"strlen": length,
		"size":   total,
		"head":   r[:5],
	}).Debug("created strings")
	return r
}

func genString(n int64) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[src.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func CreateObjects(total int64) Objects {
	r := make(Objects, total)
	for idx := range r {
		r[idx] = genObject()
	}
	logrus.WithFields(map[string]interface{}{
		"typ":  TypeObject,
		"size": total,
		"head": r[:5],
	}).Debug("created objects")
	return r
}

func genObject() Object {
	rnd := rand.New(src)
	return Object{
		I: genInt(math.MinInt64, math.MaxInt64),
		F: rnd.NormFloat64(),
		T: genInt(0, 1)%2 == 0,
		S: genString(100),
		B: []byte(genString(100))[:100], // force 100 bytes
	}
}
