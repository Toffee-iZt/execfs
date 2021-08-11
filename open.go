package workfs

import (
	"encoding/json"
	"io/fs"
	"os"
)

// Open opens file in the exec directory.
func Open(name string) (fs.File, error) {
	return AsFS().Open(name)
}

// OpenOS opens file (but as *os.File) in the exec directory.
func OpenOS(name string) (*os.File, error) {
	return os.Open(abs(name))
}

// OpenFile opens file in the exec directory with rw flag.
func OpenFile(name string) (*os.File, error) {
	return os.OpenFile(abs(name), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
}

// ReadFile reads file in the exec directory.
func ReadFile(name string) ([]byte, error) {
	return os.ReadFile(abs(name))
}

// WriteFile writes data to the file in the exec directory.
func WriteFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0666)
}

// LoadJSON opens and parses json file in the exec directory.
func LoadJSON(name string, dst interface{}) error {
	f, err := Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(dst)
}

// SaveJSON opens file in the exec directory and writes the json encoding.
func SaveJSON(name string, v interface{}) error {
	f, err := OpenFile(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}
