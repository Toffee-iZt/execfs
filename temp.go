package wfs

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var once sync.Once
var worktemp string

// OSTempDir returns the system temp directory.
func OSTempDir() string {
	return os.TempDir()
}

func tempName(suffix string) string {
	return strconv.FormatInt(time.Now().UnixNano(), 10) + suffix
}

// GetWorkTempDir returns the working temp directory.
func GetWorkTempDir(suffix string) string {
	once.Do(func() {
		if suffix == "" {
			suffix = GetExecName()
		}
		worktemp = filepath.Join(OSTempDir(), tempName(suffix))
		err := os.Mkdir(worktemp, 0700)
		if err != nil {
			panic(err)
		}
	})
	return worktemp
}

// GetTempFS returns temp filesystem.
func GetTempFS() Filesystem {
	return OpenOS(GetWorkTempDir(""))
}
