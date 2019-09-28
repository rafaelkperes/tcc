package data

import (
	"encoding/json"
	"errors"
)

// Type reference
type Type string

const (
	TypeUndefined Type = "undefined"
	TypeInt            = "int"
	TypeFloat          = "float"
	TypeString         = "string"
	TypeObject         = "object"
)

// Format for marshalling as mime string
type Format string

const (
	FormatUndefined Format = "undefined"
	FormatJSON             = "application/json"
	FormatProtobuf         = "application/x-protobuf"
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

func (t Ints) Marshal(format Format) (data []byte, typ Type, err error) {
	typ = TypeInt
	switch format {
	case FormatJSON:
		data, err = json.Marshal(t)
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
	default:
		return nil, TypeUndefined, errors.New("unknown format")
	}
}
