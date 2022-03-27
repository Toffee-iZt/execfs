package wfs

import (
	"os"
	"path/filepath"
)

var exec, execDir, execName string

func init() {
	var err error
	exec, err = os.Executable()
	if err != nil {
		panic(err)
	}
	execDir, execName = filepath.Split(exec)
}

// ExecPath returns the path to the executable file.
func ExecPath() string {
	return exec
}

// ExecDir returns the directory of the excutable file.
func ExecDir() string {
	return execDir
}

// ExecName returns the name of the excutable file.
func ExecName() string {
	return execName
}

// ExecFS returns executable filesystem.
func ExecFS() Filesystem {
	return OpenOS(execDir)
}
