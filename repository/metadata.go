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
	"sort"
	"time"
)

// PackageMetadata represents a package's metadata.
type PackageMetadata struct {
	PackageID string
	MavenData MavenData
	Versions  []PackageVersionMetadata
}

// PackageVersionMetadata represents a version in the version.
type PackageVersionMetadata struct {
	Version string
	Files   map[string]UploadedFileMetadata
	Created time.Time
}

type packageVersionCollection []PackageVersionMetadata

func (s packageVersionCollection) Len() int {
	return len(s)
}
func (s packageVersionCollection) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s packageVersionCollection) Less(i, j int) bool {
	return s[i].Created.Before(s[j].Created)
}

// UploadedFileMetadata holds some basic information about a file's storage ID and hashes.
type UploadedFileMetadata struct {
	ID   string
	SHA1 string
	MD5  string
}

// MavenData denotes a package ID for Maven packages.
type MavenData struct {
	GroupID    string
	ArtifactID string
}

// SortVersionsByCreated sorts the PackageVersionMetadata slice provided by creation date.
func SortVersionsByCreated(versions []PackageVersionMetadata) {
	sort.Sort(packageVersionCollection(versions))
}
