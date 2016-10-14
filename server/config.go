package server

import (
	"encoding/json"
	"os"
)

type Config struct {
	Listen       string
	Repositories map[string]RepoConfig
	SSL          SSLConfig
}

type RepoConfig struct {
	Metadata       MetadataConfig
	Storage        StorageConfig
	Authentication AuthenticationConfig
}

type MetadataConfig struct {
	Path string
}

type StorageConfig struct {
	Path string
}

type AuthenticationConfig struct {
	Type string
}

type SSLConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err = json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
