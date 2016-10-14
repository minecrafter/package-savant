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
