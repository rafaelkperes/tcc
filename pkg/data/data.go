package data

import (
	"encoding/json"
	"errors"

	"github.com/hamba/avro"

	"github.com/golang/protobuf/proto"
	"github.com/rafaelkperes/tcc/pkg/data/pbdata"
	"github.com/ugorji/go/codec"
)

// Type reference
type Type string

const (
	TypeUndefined Type = "undefined"
	TypeInt       Type = "int"
	TypeFloat     Type = "float"
	TypeString    Type = "string"
	TypeObject    Type = "object"
)

// Format for marshalling as mime string
type Format string

const (
	FormatUndefined Format = "undefined"
	FormatJSON      Format = "application/json"
	FormatProtobuf  Format = "application/x-protobuf"
	FormatMsgpack   Format = "application/x-msgpack"
	FormatAvro      Format = "application/x-avro"
)

type Data interface {
	// Marshal the data, returning the marshalled bytes and the type reference,
	// or an error otherwise.
	Marshal(format Format) (data []byte, typ Type, err error)
}

type Ints []int64

type Floats []float64

type Strings []string

type Objects []Object

type Object struct {
	I int64
	F float64
	T bool
	S string
	B []byte
}

var msgpHandle = codec.MsgpackHandle{}

func (t Ints) Marshal(format Format) (data []byte, typ Type, err error) {
	typ = TypeInt
	switch format {
	case FormatJSON:
		data, err = json.Marshal(t)
		return
	case FormatProtobuf:
		d := &pbdata.Ints{
			Ints: t,
		}
		data, err = proto.Marshal(d)
		return
	case FormatMsgpack:
		enc := codec.NewEncoderBytes(&data, &msgpHandle)
		err = enc.Encode(t)
		return
	case FormatAvro:
		data, err = avro.Marshal(avroInts, t)
		return
	default:
		return nil, TypeUndefined, errors.New("unknown format")
	}
}

func (t Floats) Marshal(format Format) (data []byte, typ Type, err error) {
	typ = TypeFloat
	switch format {
	case FormatJSON:
		data, err = json.Marshal(t)
		return
	case FormatProtobuf:
		d := &pbdata.Floats{
			Floats: t,
		}
		data, err = proto.Marshal(d)
		return
	case FormatMsgpack:
		enc := codec.NewEncoderBytes(&data, &msgpHandle)
		err = enc.Encode(t)
		return
	case FormatAvro:
		data, err = avro.Marshal(avroFloats, t)
		return
	default:
		return nil, TypeUndefined, errors.New("unknown format")
	}
}

func (t Strings) Marshal(format Format) (data []byte, typ Type, err error) {
	typ = TypeString
	switch format {
	case FormatJSON:
		data, err = json.Marshal(t)
		return
	case FormatProtobuf:
		d := &pbdata.Strings{
			Strings: t,
		}
		data, err = proto.Marshal(d)
		return
	case FormatMsgpack:
		enc := codec.NewEncoderBytes(&data, &msgpHandle)
		err = enc.Encode(t)
		return
	case FormatAvro:
		data, err = avro.Marshal(avroStrings, t)
		return
	default:
		return nil, TypeUndefined, errors.New("unknown format")
	}
}

func (t Objects) Marshal(format Format) (data []byte, typ Type, err error) {
	typ = TypeObject
	switch format {
	case FormatJSON:
		data, err = json.Marshal(t)
		return
	case FormatProtobuf:
		pbObjs := make([]*pbdata.Objects_Object, len(t))
		for idx, obj := range t {
			pbObjs[idx] = &pbdata.Objects_Object{
				B: obj.B,
				F: obj.F,
				I: obj.I,
				S: obj.S,
				T: obj.T,
			}
		}

		d := &pbdata.Objects{
			Objects: pbObjs,
		}
		data, err = proto.Marshal(d)
		return
	case FormatMsgpack:
		enc := codec.NewEncoderBytes(&data, &msgpHandle)
		err = enc.Encode(t)
		return
	case FormatAvro:
		data, err = avro.Marshal(avroObjects, t)
		return
	default:
		return nil, TypeUndefined, errors.New("unknown format")
	}
}
