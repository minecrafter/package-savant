package store

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/minecrafter/sage/repository"
)

// Represents a local metadata store.
type LocalMetadataStore struct {
	sync.RWMutex
	loadedData map[string]repository.PackageMetadata
	path       string
}

// Creates a new local metadata store instance.
func NewLocalMetadataStore(path string) *LocalMetadataStore {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var data map[string]repository.PackageMetadata
	if err = json.NewDecoder(file).Decode(&data); err != nil {
		panic(err)
	}
	return &LocalMetadataStore{
		loadedData: data,
		path:       path,
	}
}

// Looks up a package's metadata based on its ID.
func (ms *LocalMetadataStore) FindByID(id string) (*repository.PackageMetadata, error) {
	ms.RLock()
	data, exists := ms.loadedData[id]
	ms.RUnlock()
	if exists {
		return &data, nil
	}
	return nil, repository.ErrPackageNotFound
}

// Stores new metadata for a package.
func (ms *LocalMetadataStore) Store(metadata repository.PackageMetadata) error {
	ms.Lock()
	defer ms.Unlock()
	ms.loadedData[metadata.PackageID] = metadata

	file, err := os.OpenFile(ms.path, os.O_WRONLY, 666)
	if err != nil {
		return err
	}
	defer file.Close()

	json.NewEncoder(file).Encode(ms.loadedData)
	return nil
}
