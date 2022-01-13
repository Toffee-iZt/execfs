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

// GetExec returns the path to the executable file.
func GetExec() string {
	return exec
}

// GetExecDir returns the directory of the excutable file.
func GetExecDir() string {
	return execDir
}

// GetExecName returns the name of the excutable file.
func GetExecName() string {
	return execName
}

// GetExecFS returns executable filesystem.
func GetExecFS() Filesystem {
	return OpenOS(execDir)
}
