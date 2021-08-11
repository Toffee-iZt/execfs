// +build linux android netbsd

package workfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

// AsFS returns the file system rooted in the exec directory.
func AsFS() fs.FS {
	return GetFS("")
}

// GetFS returns the file system rooted relavite to exec directory.
func GetFS(name string) fs.FS {
	return os.DirFS(abs(name))
}

// ReadDir reads the related named directory and returns a list of directory entries.
func ReadDir(name string) ([]fs.DirEntry, error) {
	f, err := OpenOS(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.ReadDir(-1)
}

func abs(name string) string {
	return filepath.Join(dir, name)
}

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
