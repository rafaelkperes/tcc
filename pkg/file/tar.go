package file

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// AsTarball packs the given directory as a gzipped tarball, writing the
// resulting bytes in w.
func AsTarball(path string, w io.Writer) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	return tarWalk(tw, path)
}

func tarWalk(tw *tar.Writer, rootPath string) error {
	rootPath = filepath.Clean(rootPath)
	walk := func(path string, info os.FileInfo, err error) error {
		path = filepath.Clean(path)
		path = strings.TrimPrefix(path, rootPath)
		if len(path) == 0 {
			return nil
		}
		path = strings.TrimPrefix(path, "/")

		th, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		th.Name = path
		if err := tw.WriteHeader(th); err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			b, err := ioutil.ReadFile(filepath.Join(rootPath, path))
			if err != nil {
				return err
			}
			_, err = tw.Write(b)
			return err
		}
		return nil
	}
	return filepath.Walk(rootPath, walk)
}
