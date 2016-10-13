package server

import (
    "encoding/json"
    "os"
)

type Config struct {
    Repositories map[string]RepoConfig
}

type RepoConfig struct {
    Metadata MetadataConfig
    Storage StorageConfig
}

type MetadataConfig struct {
    Path string
}

type StorageConfig struct {
    Path string
}

func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    defer file.Close()
    if err != nil {
        return nil, err
    }
    
    var config Config
    if err = json.NewDecoder(file).Decode(&config); err != nil {
        return nil, err
    }

    return &config, nil
}