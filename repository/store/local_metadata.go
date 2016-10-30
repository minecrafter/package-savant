// Copyright 2016 Package Savant team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/minecrafter/package-savant/repository"
	"github.com/pkg/errors"
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
	return nil, errors.WithStack(repository.ErrPackageNotFound)
}

// Stores new metadata for a package.
func (ms *LocalMetadataStore) Store(metadata repository.PackageMetadata) error {
	ms.Lock()
	defer ms.Unlock()
	ms.loadedData[metadata.PackageID] = metadata

	file, err := os.OpenFile(ms.path, os.O_WRONLY, 666)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	json.NewEncoder(file).Encode(ms.loadedData)
	return nil
}

// Returns all package IDs.
func (ms *LocalMetadataStore) GetAllIDs() (*[]string, error) {
	keys := make([]string, len(ms.loadedData))

	i := 0
	for k := range ms.loadedData {
		keys[i] = k
		i++
	}

	return &keys, nil
}
