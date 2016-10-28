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

package repository

import (
	"errors"
	"io"
)

var (
	// ErrPackageNotFound is returned when a package is not found.
	ErrPackageNotFound = errors.New("Package not found")
)

// MetadataStore represents a storage of metadata.
type MetadataStore interface {
	FindByID(id string) (*PackageMetadata, error)
	GetAllIDs() (*[]string, error)
	Store(metadata PackageMetadata) error
}

// StorageStore represents storage for packages and other data.
type StorageStore interface {
	ReadByID(id string) (io.ReadSeeker, error)
	Write() (io.WriteCloser, string, error)
}
