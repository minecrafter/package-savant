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

import (
	"encoding/xml"
)

// Represents a maven-metadata.xml file.
type MavenMetadata struct {
	XMLName    xml.Name                `xml:"metadata"`
	GroupID    string                  `xml:"groupId"`
	ArtifactID string                  `xml:"artifactId"`
	Versioning MavenMetadataVersioning `xml:"versioning"`
}

type MavenMetadataVersioning struct {
	Latest      string                `xml:"latest"`
	Release     string                `xml:"release"`
	Versions    MavenMetadataVersions `xml:"versions"`
	LastUpdated string                `xml:"lastUpdated"`
}

type MavenMetadataVersions struct {
	Version []string `xml:"version"`
}
