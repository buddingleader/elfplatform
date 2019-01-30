package common

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// FileOrFolderExists checks if a file or folder exists
func FileOrFolderExists(fileOrFolder string) bool {
	_, err := os.Stat(fileOrFolder)
	return !os.IsNotExist(err)
}

// CreateDir creates a dir for dirPath
func CreateDir(dirPath string) error {
	dirPath = appendSuffixToPath(dirPath)

	if FileOrFolderExists(dirPath) {
		return nil
	}

	return os.MkdirAll(path.Dir(dirPath), 0766)
}

func appendSuffixToPath(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return path
}

// DirEmpty returns true if the dir at dirPath is empty
func DirEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, errors.Wrapf(err, "error opening dir [%s]", dirPath)
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, errors.Wrapf(err, "error checking if dir [%s] is empty", dirPath)
}
