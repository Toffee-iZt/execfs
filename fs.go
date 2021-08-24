package wfs

import (
	"io"
	sysfs "io/fs"
	"os"
	"path/filepath"
	"syscall"
)

// Filesystem interface.
type Filesystem interface {
	OpenFile(name string, flag int, perm sysfs.FileMode) (File, error)
	Remove(name string) error
	Rename(oldpath, newpath string) error
	Mkdir(name string, perm sysfs.FileMode) error
	Stat(name string) (sysfs.FileInfo, error)
	ReadDir(path string) ([]sysfs.DirEntry, error)
}

// File interface.
type File interface {
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
	io.Seeker
	io.Closer
	Name() string
	Stat() (sysfs.FileInfo, error)
	Sync() error
	Truncate(int64) error
}

// Open opens the named file with read flag.
func Open(fs Filesystem, name string) (File, error) {
	return fs.OpenFile(name, os.O_RDONLY, 0)
}

// Create creates or truncates the named file with rw flag.
func Create(fs Filesystem, name string) (File, error) {
	return fs.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
}

// MkdirAll creates a directory named path along with any necessary parents.
func MkdirAll(fs Filesystem, path string, perm sysfs.FileMode) error {
	if dir, err := fs.Stat(path); err == nil {
		if dir.IsDir() {
			return nil
		}
		return &sysfs.PathError{Op: "mkdir", Path: path, Err: syscall.ENOTDIR}
	}

	parent, _ := filepath.Split(path)
	if parent != "" {
		err := MkdirAll(fs, parent, perm)
		if err != nil {
			return err
		}
	}

	err := fs.Mkdir(path, perm)
	if err != nil {
		dir, err1 := fs.Stat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}

// RemoveAll removes path and any children it contains.
func RemoveAll(fs Filesystem, path string) error {
	if err := fs.Remove(path); err == nil || os.IsNotExist(err) {
		return nil
	}

	infos, err := fs.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, info := range infos {
		child := filepath.Join(path, info.Name())
		err1 := RemoveAll(fs, child)
		if err == nil {
			err = err1
		}
	}

	err1 := fs.Remove(path)
	if err1 == nil || os.IsNotExist(err1) {
		return nil
	}
	return err
}
