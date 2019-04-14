package gen

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	allIntsFile   = "ai.json"
	bigIntsFile   = "bi.json"
	smallIntsFile = "si.json"
)

// All populates dst directory with all data files.
// dst is created if it doesn't exist.
func All(dst string, size int) error {
	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	err = writeData(filepath.Join(dst, smallIntsFile), smallInts(size))
	if err != nil {
		return errors.Wrap(err, "failed to write small ints")
	}

	err = writeData(filepath.Join(dst, bigIntsFile), bigInts(size))
	if err != nil {
		return errors.Wrap(err, "failed to write small ints")
	}

	err = writeData(filepath.Join(dst, allIntsFile), allInts(size))
	if err != nil {
		return errors.Wrap(err, "failed to write small ints")
	}

	return nil
}

func writeData(name string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal data to json")
	}

	return ioutil.WriteFile(name, b, os.ModePerm)
}
