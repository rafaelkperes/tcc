package data

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"

	"github.com/hamba/avro"

	"github.com/ugorji/go/codec"

	"github.com/golang/protobuf/proto"

	"github.com/rafaelkperes/tcc/pkg/data/pbdata"

	"github.com/sirupsen/logrus"
)

// Using fixed seed to always reproduce the same scenario.
// var src = rand.NewSource(time.Now().UnixNano())
var src = rand.NewSource(88)

func Create(typ Type, total int64, opts ...Option) (Data, error) {
	options := newOptions()
	for _, o := range opts {
		o.apply(options)
	}

	switch typ {
	case TypeInt:
		return CreateInts(total, options.intMin, options.intMax), nil
	case TypeFloat:
		return CreateFloats(total), nil
	case TypeString:
		return CreateStrings(total, options.strLen), nil
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
	// protobuf needs to map between proto generated types and the actual ones
	case FormatProtobuf:
		switch typ {
		case TypeInt:
			var v pbdata.Ints
			err := proto.Unmarshal(d, &v)
			return Ints(v.GetInts()), err
		case TypeFloat:
			var v pbdata.Floats
			err := proto.Unmarshal(d, &v)
			return Floats(v.GetFloats()), err
		case TypeString:
			var v pbdata.Strings
			err := proto.Unmarshal(d, &v)
			return Strings(v.GetStrings()), err
		case TypeObject:
			var v pbdata.Objects
			err := proto.Unmarshal(d, &v)
			pbObjs := v.GetObjects()
			res := make(Objects, len(pbObjs))
			for idx, pbo := range pbObjs {
				res[idx] = Object{
					B: pbo.B,
					F: pbo.F,
					I: pbo.I,
					S: pbo.S,
					T: pbo.T,
				}
			}
			return res, err
		default:
			return nil, fmt.Errorf("unknown type: %v", typ)
		}
		// protobuf needs to map between proto generated types and the actual ones
	case FormatMsgpack:
		switch typ {
		case TypeInt:
			var v Ints
			dec := codec.NewDecoderBytes(d, &msgpHandle)
			err := dec.Decode(&v)
			return v, err
		case TypeFloat:
			var v Floats
			dec := codec.NewDecoderBytes(d, &msgpHandle)
			err := dec.Decode(&v)
			return v, err
		case TypeString:
			var v Strings
			dec := codec.NewDecoderBytes(d, &msgpHandle)
			err := dec.Decode(&v)
			return v, err
		case TypeObject:
			var v Objects
			dec := codec.NewDecoderBytes(d, &msgpHandle)
			err := dec.Decode(&v)
			return v, err
		default:
			return nil, fmt.Errorf("unknown type: %v", typ)
		}
	case FormatAvro:
		switch typ {
		case TypeInt:
			var v Ints
			err := avro.Unmarshal(avroInts, d, &v)
			return v, err
		case TypeFloat:
			var v Floats
			err := avro.Unmarshal(avroFloats, d, &v)
			return v, err
		case TypeString:
			var v Strings
			err := avro.Unmarshal(avroStrings, d, &v)
			return v, err
		case TypeObject:
			var v Objects
			err := avro.Unmarshal(avroObjects, d, &v)
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

	var peak int64 = 5
	if total < 5 {
		peak = total
	}
	logrus.WithFields(map[string]interface{}{
		"typ":    TypeInt,
		"size":   total,
		"intmin": min,
		"intmax": max,
		"head":   r[:peak],
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

	var peak int64 = 5
	if total < 5 {
		peak = total
	}
	logrus.WithFields(map[string]interface{}{
		"typ":  TypeInt,
		"size": total,
		"head": r[:peak],
	}).Debug("created floats")
	return r
}

func CreateStrings(total, length int64) Strings {
	r := make(Strings, total)
	for idx := range r {
		r[idx] = genString(length)
	}

	var peak int64 = 5
	if total < 5 {
		peak = total
	}
	logrus.WithFields(map[string]interface{}{
		"typ":    TypeString,
		"strlen": length,
		"size":   total,
		"head":   r[:peak],
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

	var peak int64 = 5
	if total < 5 {
		peak = total
	}
	logrus.WithFields(map[string]interface{}{
		"typ":  TypeObject,
		"size": total,
		"head": r[:peak],
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
