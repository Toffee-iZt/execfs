package wfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

var _ Filesystem = &osfs{}

// OpenOS opens os file as Filesystem.
func OpenOS(dir string) Filesystem {
	return &osfs{dir: dir}
}

type osfs struct {
	dir string
}

func (f *osfs) abs(name string) string {
	return filepath.Join(f.dir, name)
}

func (f *osfs) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	return os.OpenFile(f.abs(name), flag, perm)
}

func (f *osfs) Remove(name string) error {
	return os.Remove(f.abs(name))
}

func (f *osfs) Rename(oldpath, newpath string) error {
	return os.Rename(f.abs(oldpath), f.abs(newpath))
}

func (f *osfs) Mkdir(name string, perm fs.FileMode) error {
	return os.Mkdir(f.abs(name), perm)
}

func (f *osfs) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(f.abs(name))
}

func (f *osfs) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(f.abs(name))
}
