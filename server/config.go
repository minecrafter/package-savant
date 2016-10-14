package server

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"golang.org/x/crypto/bcrypt"
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
	Logins map[string]string
}

type SSLConfig struct {
	Enabled bool
	CertFile string
	KeyFile string
}

func (c AuthenticationConfig) CheckPassword(username, password string) (bool, error) {
	b64password, ok := c.Logins[username]
	if !ok {
		return false, nil
	}

	decodedPassword, err := base64.StdEncoding.DecodeString(b64password)
	if err != nil {
		return false, err
	}

	if err = bcrypt.CompareHashAndPassword(decodedPassword, []byte(password)); err != nil {
		return false, err
	}
	return true, nil
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
