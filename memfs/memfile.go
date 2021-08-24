package memfs

import (
	"errors"
	"io"
	"io/fs"
	"time"
)

// MemFile is an in-memory file descriptor.
type MemFile struct {
	f      *file
	name   string
	offset int64
	closed bool
}

func (fd *MemFile) isClosed() bool {
	return fd.closed || fd.f == nil
}

func (fd *MemFile) Read(b []byte) (int, error) {
	if fd.isClosed() {
		return 0, fs.ErrClosed
	}
	n, err := fd.f.read(b, fd.offset)
	if err != nil {
		return 0, err
	}
	fd.offset += int64(n)
	return n, nil
}

func (fd *MemFile) ReadAt(b []byte, off int64) (int, error) {
	if fd.isClosed() {
		return 0, fs.ErrClosed
	}
	n, err := fd.f.read(b, off)
	if err != nil {
		return 0, err
	}
	if n < len(b) {
		err = io.EOF
	}
	return n, err
}

func (fd *MemFile) Write(b []byte) (int, error) {
	if fd.isClosed() {
		return 0, fs.ErrClosed
	}
	n, err := fd.f.write(b, fd.offset)
	if err != nil {
		return 0, err
	}
	fd.offset += int64(n)
	return n, nil
}

func (fd *MemFile) Seek(offset int64, whence int) (int64, error) {
	if fd.isClosed() {
		return 0, fs.ErrClosed
	}
	switch whence {
	default:
	case 1:
		offset += fd.offset
	case 2:
		offset += fd.f.size()
	}
	if offset < 0 {
		return 0, errors.New("seek: offset before file begin")
	}
	fd.offset = offset
	return offset, nil
}

// Close returns nil because mem file cannot be closed.
func (fd *MemFile) Close() error {
	if fd.isClosed() {
		return fs.ErrClosed
	}
	fd.closed = true
	return nil
}

// Name returns name of the mem file.
func (fd *MemFile) Name() string {
	return fd.name
}

// Stat implements fs.File.Stat.
func (fd *MemFile) Stat() (fs.FileInfo, error) {
	if fd.isClosed() {
		return nil, fs.ErrClosed
	}
	return &fileInfo{
		name:    fd.f.name,
		size:    fd.f.size(),
		modTime: fd.f.modTime,
		mode:    fd.f.perm,
	}, nil
}

// Name returns name of the mem file.
func (fd *MemFile) Truncate(size int64) error {
	if fd.isClosed() {
		return fs.ErrClosed
	}
	return fd.f.trunc(size)
}

var ErrNotRegular = errors.New("file is not regular")

type file struct {
	name     string
	perm     fs.FileMode
	data     []byte
	children []*file
	modTime  time.Time
}

func (f *file) isDir() bool {
	return f.perm.IsDir()
}

func (f *file) size() int64 {
	return int64(len(f.data))
}

func (f *file) grow(n int64) (int64, error) {
	l := int64(len(f.data))
	c := int64(cap(f.data))
	if n > c-l {
		newdata := make([]byte, c*2+n)
		l = int64(copy(newdata, f.data))
		f.data = newdata[:l+n]
	}
	return l, nil
}

func (f *file) read(b []byte, off int64) (int, error) {
	if f.isDir() {
		return 0, ErrNotRegular
	}
	if f.size() <= off {
		return 0, io.EOF
	}
	return copy(b, f.data[off:]), nil
}

func (f *file) write(b []byte, off int64) (int, error) {
	if f.isDir() {
		return 0, ErrNotRegular
	}
	n := int64(len(b))
	f.grow(n + off - f.size())
	return copy(f.data[off:off+n], b), nil
}

func (f *file) trunc(size int64) error {
	if f.isDir() {
		return ErrNotRegular
	}
	if size < 0 {
		return errors.New("truncate: size must be non-negative")
	}
	fsize := f.size()
	if size > fsize {
		f.grow(fsize - size)
	} else {
		f.data = f.data[:size]
	}
	return nil
}

type fileInfo struct {
	name    string
	size    int64
	modTime time.Time
	mode    fs.FileMode
}

// Name returns base name of the file.
func (i *fileInfo) Name() string { return i.name }

// Size returns file length in bytes.
func (i *fileInfo) Size() int64 { return i.size }

// Mode returns file mode bits.
func (i *fileInfo) Mode() fs.FileMode { return i.mode }

// ModTime returns modification time.
func (i *fileInfo) ModTime() time.Time { return i.modTime }

// IsDir is abbreviation for Mode().IsDir().
func (i *fileInfo) IsDir() bool { return i.mode&fs.ModeDir > 0 }

// Sys returns underlying data source (for memory file returns nil).
func (i *fileInfo) Sys() interface{} { return nil }
