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
