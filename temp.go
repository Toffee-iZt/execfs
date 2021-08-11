package workfs

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// CreateTemp creates temp file.
func CreateTemp(suffix string) (*os.File, error) {
	path, err := tempPath(suffix)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(tempabs(path), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
}

// CreateTempDir creates temp directory.
func createTempDir(suffix string) (string, error) {
	path, err := tempPath(suffix)
	if err != nil {
		return "", err
	}
	return path, os.Mkdir(tempabs(path), 0700)
}

var tempdir, worktemp string

func init() {
	tempdir = os.TempDir()
	var err error
	worktemp, err = createTempDir(GetExecName())
	if err != nil {
		panic(err)
	}
}

// GetTempDir returns the system temp directory.
func GetTempDir() string {
	return dir
}

// GetWorkTemp returns the working temp directory.
func GetWorkTemp() string {
	return name
}

func tempabs(name string) string {
	return filepath.Join(worktemp, name)
}

var errHasSeparator = errors.New("suffix contains path separator")

func tempPath(suffix string) (string, error) {
	for i := 0; i < len(suffix); i++ {
		if os.IsPathSeparator(suffix[i]) {
			return "", &os.PathError{Op: "createtemp", Path: suffix, Err: errHasSeparator}
		}
	}
	return tempabs(strconv.FormatInt(time.Now().UnixNano(), 10) + suffix), nil
}
