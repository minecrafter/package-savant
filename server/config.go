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
