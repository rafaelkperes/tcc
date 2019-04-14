package gen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratedFiles(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("./", t.Name())
	r.NoError(err)
	defer os.RemoveAll(dir)

	// generate files
	r.NoError(All(dir, 1))

	fis, err := ioutil.ReadDir(dir)
	r.NoError(err)
	t.Logf("generated %d files", len(fis))

	names := make([]string, len(fis))
	for idx, fi := range fis {
		names[idx] = fi.Name()
	}

	r.Contains(names, filepath.Base(smallIntsFile))
	r.Contains(names, filepath.Base(bigIntsFile))
	r.Contains(names, filepath.Base(allIntsFile))
}

func TestWriteData(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("./", t.Name())
	r.NoError(err)
	defer os.RemoveAll(dir)

	name := filepath.Join(dir, "f")
	r.NoError(writeData(name, []string{"foo", "bar"}))

	b, err := ioutil.ReadFile(name)
	r.NoError(err)
	// expect data as JSON
	r.Equal(`["foo","bar"]`, string(b))
}
