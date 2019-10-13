package file

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTar(t *testing.T) {
	r := require.New(t)
	var b []byte
	buff := bytes.NewBuffer(b)

	r.NoError(AsTarball("./test", buff))

	want := map[string][]byte{
		"foo.txt":        []byte("foo"),
		"subdir":         nil,
		"subdir/bar.txt": []byte("bar"),
	}

	got := map[string][]byte{}
	gr, err := gzip.NewReader(buff)
	r.NoError(err)
	tr := tar.NewReader(gr)
	for th, err := tr.Next(); err != io.EOF; th, err = tr.Next() {
		r.NoError(err)

		if !th.FileInfo().Mode().IsRegular() {
			got[th.Name] = nil
			continue
		}

		buff := bytes.NewBuffer([]byte{})
		_, err := io.Copy(buff, tr)
		r.NoError(err)
		got[th.Name] = buff.Bytes()
	}

	r.Equal(want, got)
}
