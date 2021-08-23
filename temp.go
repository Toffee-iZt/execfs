package wfs

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var tempdir = os.TempDir()
var worktemp string
var wconce sync.Once

// GetTempDir returns the system temp directory.
func GetTempDir() string {
	return tempdir
}

// GetWorkTemp returns the working temp directory.
func GetWorkTemp() string {
	return worktemp
}

// CreateTemp creates temp file.
func CreateTemp(suffix string) (*os.File, error) {
	path, err := tempPath(suffix)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
}

// CreateTempDir creates temp directory.
func CreateTempDir(suffix string) (string, error) {
	path, err := tempPath(suffix)
	if err != nil {
		return "", err
	}
	return path, os.Mkdir(path, 0700)
}

func makeWorkTempDir() {
	wconce.Do(func() {
		worktemp = filepath.Join(tempdir, tempName(GetExecName()))
		err := os.Mkdir(worktemp, 0700)
		if err != nil {
			panic(err)
		}
	})
}

var errHasSeparator = errors.New("suffix contains path separator")

func tempPath(suffix string) (string, error) {
	makeWorkTempDir()
	for i := 0; i < len(suffix); i++ {
		if os.IsPathSeparator(suffix[i]) {
			return "", &os.PathError{Op: "createtemp", Path: suffix, Err: errHasSeparator}
		}
	}
	return filepath.Join(worktemp, tempName(suffix)), nil
}

func tempName(suffix string) string {
	return strconv.FormatInt(time.Now().UnixNano(), 10) + suffix
}
