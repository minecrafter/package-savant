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

package maven

import "github.com/minecrafter/sage/repository"

const (
	mavenDateFormat = "20060102150405"
)

// CreateMavenMetadata converts a Sage PackageMetadata into a maven-metadata.xml file.
func CreateMavenMetadata(metadata repository.PackageMetadata) MavenMetadata {
	m := MavenMetadata{
		GroupID:    metadata.MavenData.GroupID,
		ArtifactID: metadata.MavenData.ArtifactID,
	}

	// Copy and sort versions. While we could have sorted versions beforehand, we would only gain minimal benefits.
	versions := make([]repository.PackageVersionMetadata, len(metadata.Versions))
	copy(versions, metadata.Versions)
	repository.SortVersionsByCreated(versions)

	// Stringify the sorted versions.
	sortedVersions := make([]string, len(metadata.Versions))
	for i, version := range versions {
		sortedVersions[i] = version.Version
	}
	m.Versioning = MavenMetadataVersioning{
		Latest:      sortedVersions[len(sortedVersions)-1],
		Versions:    MavenMetadataVersions{sortedVersions},
		LastUpdated: versions[len(versions)-1].Created.Format(mavenDateFormat),
	}
	return m
}
