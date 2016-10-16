package maven

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/minecrafter/sage/repository"
	"github.com/minecrafter/sage/util"
)

type MavenRetrieveHandler struct {
	root          string
	metadataStore repository.MetadataStore
	storageStore  repository.StorageStore
}

func NewMavenRetrieveHandler(path string, metadataStore repository.MetadataStore, storageStore repository.StorageStore) MavenRetrieveHandler {
	return MavenRetrieveHandler{
		root:          path,
		metadataStore: metadataStore,
		storageStore:  storageStore,
	}
}

func (h *MavenRetrieveHandler) GetMavenFile(w http.ResponseWriter, r *http.Request) {
	// Example URLs:
	// http://repo.maven.apache.org/maven2/com/aerospike/aerospike-client/maven-metadata.xml
	// http://repo.maven.apache.org/maven2/com/aerospike/aerospike-client/3.3.0/aerospike-client-3.3.0.jar
	// http://repo.maven.apache.org/maven2/com/aerospike/aerospike-client/3.3.0/aerospike-client-3.3.0.pom
	path := r.URL.EscapedPath()
	if !strings.HasPrefix(path, h.root) {
		// Can't handle this request
		util.Do404(w)
		return
	}

	subtype := strings.TrimPrefix(path, h.root)

	if strings.HasSuffix(subtype, "maven-metadata.xml") {
		h.serveMavenMetadata(w, subtype)
	} else if strings.HasSuffix(subtype, ".pom") || strings.HasSuffix(subtype, ".jar") {
		h.serveMavenFile(w, r, subtype)
	} else if strings.HasSuffix(subtype, ".pom.md5") || strings.HasSuffix(subtype, ".pom.sha1") || strings.HasSuffix(subtype, ".jar.md5") || strings.HasSuffix(subtype, ".jar.sha1") {
		h.serveMavenHash(w, subtype)
	} else if subtype == "/api/packages.json" {
		h.servePackageListing(w, r)
	} else if strings.HasPrefix(subtype, "/api/packages/") && strings.HasSuffix(subtype, ".json") {
		h.serveVersionListing(w, r, subtype)
	} else {
		util.Do404(w)
	}
}

func (h *MavenRetrieveHandler) servePackageListing(w http.ResponseWriter, r *http.Request) {
	ids, err := h.metadataStore.GetAllIDs()
	if err != nil {
		log.Printf("Unable to get package ID list: %s", err.Error())
		util.DoSpecificError(w, err)
		return
	}

	out := struct {
		Packages []string `json:"packages"`
	}{
		Packages: *ids,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&out)
}

func (h *MavenRetrieveHandler) serveVersionListing(w http.ResponseWriter, r *http.Request, subtype string) {
	stripped := strings.TrimPrefix(strings.TrimSuffix(subtype, ".json"), "/api/packages/")
	pkg, err := h.metadataStore.FindByID(stripped)
	if err != nil {
		if err == repository.ErrPackageNotFound {
			// package not found
			util.Do404(w)
		} else {
			log.Printf("Unable to lookup package %s: %s", stripped, err.Error())
			util.DoSpecificError(w, err)
		}
		return
	}

	// Copy and sort versions. While we could have sorted versions beforehand, we would only gain minimal benefits.
	versions := make([]repository.PackageVersionMetadata, len(pkg.Versions))
	copy(versions, pkg.Versions)
	repository.SortVersionsByCreated(versions)

	// Stringify the sorted versions.
	sortedVersions := make([]string, len(pkg.Versions))
	for i, version := range versions {
		sortedVersions[i] = version.Version
	}

	out := struct {
		Versions []string `json:"versions"`
	}{
		Versions: sortedVersions,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&out)
}

func getPackageID(path string, inPackage bool) string {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	parts := strings.Split(path, "/")
	if inPackage {
		return strings.Join(parts[:len(parts)-3], ".") + ":" + parts[len(parts)-3]
	}
	return strings.Join(parts[:len(parts)-2], ".") + ":" + parts[len(parts)-2]
}

func getPackageMavenData(path string, inPackage bool) repository.MavenData {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	parts := strings.Split(path, "/")
	if inPackage {
		ourParts := parts[:len(parts)-2]
		return repository.MavenData{
			GroupID:    strings.Join(ourParts[:len(ourParts)-1], "."),
			ArtifactID: ourParts[len(ourParts)-1],
		}
	}
	return repository.MavenData{
		GroupID:    strings.Join(parts[:len(parts)-1], "."),
		ArtifactID: parts[len(parts)-1],
	}
}

func getPackageVersion(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-2]
}

func (h *MavenRetrieveHandler) serveMavenMetadata(w http.ResponseWriter, path string) {
	// get package ID
	id := getPackageID(path, false)

	// get package metadata
	metadata, err := h.metadataStore.FindByID(id)
	if err != nil {
		if err == repository.ErrPackageNotFound {
			// package not found
			util.Do404(w)
		} else {
			log.Printf("Unable to lookup package %s: %s", id, err.Error())
			util.DoSpecificError(w, err)
		}
		return
	}

	mavenData := CreateMavenMetadata(*metadata)
	w.Header().Add("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(mavenData)
}

func (h *MavenRetrieveHandler) serveMavenHash(w http.ResponseWriter, path string) {
	// get package ID and version
	id := getPackageID(path, true)
	version := getPackageVersion(path)

	// get package metadata
	metadata, err := h.metadataStore.FindByID(id)
	if err != nil {
		if err == repository.ErrPackageNotFound {
			// package not found
			util.Do404(w)
		} else {
			log.Printf("Unable to lookup package %s: %s", id, err.Error())
			util.DoSpecificError(w, err)
		}
		return
	}

	// get specific file requested
	for _, versionMetadata := range metadata.Versions {
		if versionMetadata.Version == version {
			// get file
			fileName := path[strings.LastIndex(path, "/"):strings.LastIndex(path, ".")]
			data, exists := versionMetadata.Files[fileName]
			if !exists {
				break
			}

			w.Header().Add("Content-Type", "text/plain")
			if strings.HasSuffix(path, "sha1") {
				w.Write([]byte(data.SHA1))
			} else if strings.HasSuffix(path, "md5") {
				w.Write([]byte(data.MD5))
			}
			return
		}
	}

	// otherwise, not found
	util.Do404(w)
}

func (h *MavenRetrieveHandler) serveMavenFile(w http.ResponseWriter, r *http.Request, path string) {
	// get package ID and version
	id := getPackageID(path, true)
	version := getPackageVersion(path)

	// get package metadata
	metadata, err := h.metadataStore.FindByID(id)
	if err != nil {
		if err == repository.ErrPackageNotFound {
			// package not found
			util.Do404(w)
		} else {
			log.Printf("Unable to lookup package %s: %s", id, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			util.DoSpecificError(w, err)
		}
		return
	}

	// get specific file requested
	for _, versionMetadata := range metadata.Versions {
		if versionMetadata.Version == version {
			// get file
			fileName := path[strings.LastIndex(path, "/"):]
			data, exists := versionMetadata.Files[fileName]
			if !exists {
				break
			}
			reader, err := h.storageStore.ReadByID(data.ID)
			if err != nil {
				util.DoSpecificError(w, err)
			} else {
				closer, ok := reader.(io.Closer)
				if ok {
					defer closer.Close()
				}
				http.ServeContent(w, r, fileName, versionMetadata.Created, reader)
			}
			return
		}
	}

	// otherwise, not found
	util.Do404(w)
}

func (h *MavenRetrieveHandler) PutMavenFile(w http.ResponseWriter, r *http.Request, path string) {
	// TODO: Authentication
	/*_, _, ok := r.BasicAuth()
	if !ok {
		w.Header.Add("WWW-Authenticate", "Basic realm=\"Sage\"")
		w.WriteHeader(http.StatusForbidden)
		return
	}*/

	// TODO: Properly handle input. Check if we're uploading a real JAR or POM file for instance.

	// get package ID and version
	id := getPackageID(path, true)
	version := getPackageVersion(path)

	// Copy body to new storage file. Use io.TeeReader to allow us to calculate hashes at the same time.
	writer, sid, err := h.storageStore.Write()
	if err != nil {
		log.Printf("Unable to create content: %s", err.Error())
		util.DoSpecificError(w, err)
		return
	}
	defer writer.Close()
	sha1Sum := sha1.New()
	md5Sum := md5.New()
	teeReader := io.TeeReader(r.Body, io.MultiWriter(sha1Sum, md5Sum))
	if _, err := io.Copy(writer, teeReader); err != nil {
		log.Printf("Unable to create content: %s", err.Error())
		util.DoSpecificError(w, err)
		return
	}

	fileMetadata := repository.UploadedFileMetadata{
		ID:   sid,
		SHA1: hex.EncodeToString(sha1Sum.Sum(nil)),
		MD5:  hex.EncodeToString(md5Sum.Sum(nil)),
	}
	fileName := path[strings.LastIndex(path, "/"):]

	// create version and package if needed
	metadata, err := h.metadataStore.FindByID(id)
	if err != nil {
		if err == repository.ErrPackageNotFound {
			// package not found, create
			files := make(map[string]repository.UploadedFileMetadata)
			files[fileName] = fileMetadata
			metadata = &repository.PackageMetadata{
				PackageID: id,
				MavenData: getPackageMavenData(path, true),
				Versions: []repository.PackageVersionMetadata{
					repository.PackageVersionMetadata{
						Version: version,
						Files:   files,
						Created: time.Now(),
					},
				},
			}
		} else {
			log.Printf("Unable to lookup package %s: %s", id, err.Error())
			util.DoSpecificError(w, err)
			return
		}
	} else {
		updated := false
		for _, value := range metadata.Versions {
			if value.Version == version {
				value.Files[fileName] = fileMetadata
				updated = true
				break
			}
		}

		if !updated {
			// otherwise, create and save
			files := make(map[string]repository.UploadedFileMetadata)
			files[fileName] = fileMetadata
			metadata.Versions = append(metadata.Versions, repository.PackageVersionMetadata{
				Version: version,
				Files:   files,
				Created: time.Now(),
			})
		}
	}

	if err = h.metadataStore.Store(*metadata); err != nil {
		log.Printf("Unable to create version: %s", err.Error())
		util.DoSpecificError(w, err)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}
