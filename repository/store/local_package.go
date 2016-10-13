package store

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	errPathInvalid = errors.New("Path is not valid")
)

type LocalPackageStorage struct {
	basePath string
}

func NewLocalPackageStore(path string) *LocalPackageStorage {
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &LocalPackageStorage{
		basePath: path,
	}
}

func (s *LocalPackageStorage) ReadByID(id string) (io.ReadSeeker, error) {
	realPath := filepath.Join(s.basePath, id)
	if strings.HasPrefix(realPath, s.basePath) {
		return os.Open(realPath)
	}
	return nil, errPathInvalid
}

func (s *LocalPackageStorage) Write() (io.WriteCloser, string, error) {
	// ID is 32 bytes of cryptographically random content
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, "", err
	}

	id := hex.EncodeToString(randomBytes)
	realPath := filepath.Join(s.basePath, id)
	file, err := os.Create(realPath)
	if err != nil {
		return nil, "", err
	}
	return file, id, nil
}
