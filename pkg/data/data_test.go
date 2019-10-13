package data

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataMarshalling(t *testing.T) {
	formats := []Format{
		FormatJSON,
		FormatProtobuf,
		FormatMsgpack,
		FormatAvro,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			format := format
			t.Run("Ints", func(t *testing.T) {
				r := require.New(t)
				etyp := TypeInt

				is, err := Create(etyp, 5, WithInts(0, 10))
				r.NoError(err)

				m, typ, err := is.Marshal(format)
				r.NoError(err)
				r.Equal(etyp, typ)

				d, err := Unmarshal(m, etyp, format)
				r.NoError(err)

				t.Logf("original: %+v", is)
				t.Logf("unmarshalled: %+v", d)
				r.Equal(is, d)
			})

			t.Run("Floats", func(t *testing.T) {
				r := require.New(t)
				etyp := TypeFloat

				is, err := Create(etyp, 5)
				r.NoError(err)

				m, typ, err := is.Marshal(format)
				r.NoError(err)
				r.Equal(etyp, typ)

				d, err := Unmarshal(m, etyp, format)
				r.NoError(err)

				t.Logf("original: %+v", is)
				t.Logf("unmarshalled: %+v", d)
				r.Equal(is, d)
			})

			t.Run("Strings", func(t *testing.T) {
				r := require.New(t)
				etyp := TypeString

				is, err := Create(etyp, 5, WithStringLength(10))
				r.NoError(err)

				m, typ, err := is.Marshal(format)
				r.NoError(err)
				r.Equal(etyp, typ)

				d, err := Unmarshal(m, etyp, format)
				r.NoError(err)

				t.Logf("original: %+v", is)
				t.Logf("unmarshalled: %+v", d)
				r.Equal(is, d)
			})

			t.Run("Objects", func(t *testing.T) {
				r := require.New(t)
				etyp := TypeObject

				is, err := Create(etyp, 5, WithInts(0, 10), WithStringLength(10))
				r.NoError(err)

				m, typ, err := is.Marshal(format)
				r.NoError(err)
				r.Equal(etyp, typ)

				d, err := Unmarshal(m, etyp, format)
				r.NoError(err)

				t.Logf("original: %+v", is)
				t.Logf("unmarshalled: %+v", d)
				r.Equal(is, d)
			})
		})
	}
}

func TestDataPrint(t *testing.T) {
	formats := []Format{
		FormatJSON,
		FormatProtobuf,
		FormatMsgpack,
		FormatAvro,
	}

	r := require.New(t)
	d := CreateInts(2, 0, math.MaxInt8)

	res, _, _ := d.Marshal(FormatJSON)
	t.Logf("data: %s", string(res))

	for _, format := range formats {
		res, _, err := d.Marshal(format)
		r.NoError(err)
		t.Logf("%s (%d): \t%+v", format, len(res), res)
	}
}
