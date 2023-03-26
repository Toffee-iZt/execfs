package wfs

import (
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var once sync.Once
var worktemp string

// OSTempDir returns the system temp directory.
func OSTempDir() string {
	return os.TempDir()
}

// OSTempFS returns the system temp directory as Filesystem.
func OSTempFS() Filesystem {
	return OpenOS(OSTempDir())
}

// WorkTempDir creates and returns the working temp directory.
func WorkTempDir(prefix string) string {
	once.Do(func() {
		if prefix == "" {
			prefix = ExecName()
		}
		f, err := ExecFS().Open(ExecName())
		if err != nil {
			panic(err)
		}

		h := crc32.NewIEEE()
		if _, err := io.Copy(h, f); err != nil {
			panic(err)
		}
		f.Close()

		path := filepath.Join(OSTempDir(), prefix+string(h.Sum(nil)))
		err = os.Mkdir(path, 0700)
		if !os.IsExist(err) {
			panic(err)
		}
	})
	return worktemp
}

// TempFS returns the temp filesystem for current executable.
func TempFS() Filesystem {
	return OpenOS(WorkTempDir(""))
}
