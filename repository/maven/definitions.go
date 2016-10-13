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
