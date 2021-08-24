package wfs

import (
	"encoding/json"
	"io/fs"
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

func abs(name string) string {
	return filepath.Join(dir, name)
}

// AsFS returns the file system rooted in the exec directory.
func AsFS() fs.FS {
	return os.DirFS(dir)
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

// ExecOpen opens file in the exec directory.
func ExecOpen(name string) (fs.File, error) {
	return OpenOS(name)
}

// OpenOS opens file (but as *os.File) in the exec directory.
func OpenOS(name string) (*os.File, error) {
	return os.Open(abs(name))
}

// OpenFile opens file in the exec directory with rw flag.
func OpenFile(name string, trunc bool) (*os.File, error) {
	flag := os.O_RDWR | os.O_CREATE
	if trunc {
		flag |= os.O_TRUNC
	}
	return os.OpenFile(abs(name), flag, 0666)
}

// ReadFile reads file in the exec directory.
func ReadFile(name string) ([]byte, error) {
	return os.ReadFile(abs(name))
}

// WriteFile writes data to the file in the exec directory.
func WriteFile(name string, data []byte) error {
	return os.WriteFile(abs(name), data, 0666)
}

// LoadJSON opens and parses json file in the exec directory.
func LoadJSON(name string, dst interface{}) error {
	f, err := ExecOpen(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(dst)
}

// SaveJSON opens file in the exec directory and writes the json encoding.
func SaveJSON(name string, v interface{}) error {
	f, err := OpenFile(name, true)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}
