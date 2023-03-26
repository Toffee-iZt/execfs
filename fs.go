package wfs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Filesystem interface.
type Filesystem interface {
	Open(name string) (File, error)
	OpenFile(name string, flag int, perm fs.FileMode) (File, error)
	ChangeDir(name string) Filesystem
	Create(name string) (File, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldpath, newpath string) error
	Mkdir(name string, perm fs.FileMode) error
	MkdirAll(path string, perm fs.FileMode) error
	Stat(name string) (fs.FileInfo, error)
	ReadDir(path string) ([]fs.DirEntry, error)
}

// File interface.
type File interface {
	io.Reader
	io.ReaderAt
	io.Writer
	io.Seeker
	io.Closer
	Name() string
	Stat() (fs.FileInfo, error)
	Truncate(int64) error
}

// SplitPath cleans and split path.
func SplitPath(path string) []string {
	path = filepath.Clean(path)
	parts := strings.Split(path, string(os.PathSeparator))
	for i := 0; i < len(parts); i++ {
		p := parts[i]
		if p == "" {
			copy(parts[i:], parts[i+1:])
			parts = parts[:len(parts)-1]
		}
	}
	return parts
}
