package store

import (
	"crypto/sha256"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileStore struct {
	dir string
}

func FileStoreFactory() (*Store, error) {
	home := os.Getenv("HOME")
	sd := filepath.Join(home, ".local/share/hitstore/")
	os.MkdirAll(sd, os.ModePerm)

	st := Store(FileStore{sd})
	return &st, nil
}

func (fs FileStore) Insert(key string, el Element) error {
	fn := filepath.Join(fs.dir, filename(key))
	
	err := os.Remove(fn)
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		return err
	}

	file, err := os.Create(fn)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = file.WriteString(strings.Join([]string{el.Value, el.Command}, "\n"))
	
	return nil
}

func (fs FileStore) Lookup(key string) (Element, error) {
	file, err := os.Open(filepath.Join(fs.dir, filename(key)))
	defer file.Close()
	if err != nil {
		return Element{}, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return Element{}, err
	}
	fields := strings.Split(string(b), "\n")

	return Element{Value: fields[0], Command: fields[1]}, nil
}

func (fs FileStore) Remove(key string) bool {
	return false
}

var digits = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
func filename(key string) string {
	hash := sha256.Sum256([]byte(key))
	
	b := strings.Builder{}
	for i := 0; i < len(hash); i++ {
		v := hash[i]
		b.WriteByte(digits[(v & 0xf0)>>4])
		b.WriteByte(digits[v & 0x0f])
	}

	return b.String()
}
