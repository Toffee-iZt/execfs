package workfs

import (
	"io"
	"io/fs"
	"sync"
	"time"
)

var atOncePool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 1024)
	},
}

// OpenAtOnce open FileAtOnce.
func OpenAtOnce(name string) *FileAtOnce {
	return &FileAtOnce{
		name: name,
		data: atOncePool.Get().([]byte),
	}
}

// FileAtOnce is a in-memory file that is deleted when you close it
type FileAtOnce struct {
	name string
	data []byte
	off  int
}

// Stat implements fs.File.Stat.
func (f *FileAtOnce) Stat() (fs.FileInfo, error) {
	return f, nil
}

func (f *FileAtOnce) Read(b []byte) (int, error) {
	if len(f.data) <= f.off {
		f.data = f.data[:0]
		return 0, io.EOF
	}
	n := copy(b, f.data[f.off:])
	f.off += n
	return n, nil
}

// Close clears in memory file and returns bytes to the pool.
func (f *FileAtOnce) Close() error {
	f.name = ""
	f.data = f.data[:0]
	atOncePool.Put(f.data)
	return nil
}

func (f *FileAtOnce) Write(b []byte) (int, error) {
	l := len(f.data)
	c := cap(f.data)
	n := len(b)
	m := l
	if n > c-l {
		d := f.data[f.off:]
		// can we just delete the readed data?
		if n > c-f.off {
			// not enough space anywhere, we need to allocate.
			f.data = make([]byte, c*2+n)
		}
		m = copy(f.data, d)
		f.off = 0
	}
	f.data = f.data[:m+n]
	return copy(f.data[m:m+n], b), nil
}

func (f *FileAtOnce) Name() string       { return f.name }
func (f *FileAtOnce) Size() int64        { return int64(len(f.data)) }
func (f *FileAtOnce) Mode() fs.FileMode  { return 0666 }
func (f *FileAtOnce) ModTime() time.Time { return time.Now() }
func (f *FileAtOnce) IsDir() bool        { return false }
func (f *FileAtOnce) Sys() interface{}   { return nil }
