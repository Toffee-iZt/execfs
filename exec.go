// +build linux android netbsd

package execfs

import (
	"os"
	"path/filepath"
)

var ex, dir, name string

func init() {
	var err error
	ex, err = os.Executable()
	if err != nil {
		panic(err)
	}
	dir, name = filepath.Split(ex)
}

// GetExec returns the path to the executable file.
func GetExec() string {
	return ex
}

// GetExecDir returns the directory of the excutable file.
func GetExecDir() string {
	return dir
}

// GetExecName returns the name of the excutable file.
func GetExecName() string {
	return name
}
