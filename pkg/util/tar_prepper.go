/*
Copyright 2017 Google, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/containers/image/docker/tarfile"
	"github.com/docker/docker/client"
	"github.com/golang/glog"
)

type TarPrepper struct {
	Source string
	Client *client.Client
}

func (p TarPrepper) Name() string {
	return "Tar Archive"
}

func (p TarPrepper) GetSource() string {
	return p.Source
}

func (p TarPrepper) GetImage() (Image, error) {
	return getImage(p)
}

func (p TarPrepper) GetFileSystem() (string, error) {
	return getImageFromTar(p.Source)
}

func (p TarPrepper) GetConfig() (ConfigSchema, error) {
	tempDir := strings.TrimSuffix(p.Source, filepath.Ext(p.Source)) + "-config"
	defer os.RemoveAll(tempDir)
	err := UnTar(p.Source, tempDir)
	if err != nil {
		return ConfigSchema{}, err
	}

	var config ConfigSchema
	// First open the manifest, then find the referenced config.
	manifestPath := filepath.Join(tempDir, "manifest.json")
	contents, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return ConfigSchema{}, err
	}

	manifests := []tarfile.ManifestItem{}
	if err := json.Unmarshal(contents, &manifests); err != nil {
		return ConfigSchema{}, err
	}

	if len(manifests) != 1 {
		return ConfigSchema{}, errors.New("specified tar file contains multiple images")
	}

	cfgFilename := filepath.Join(tempDir, manifests[0].Config)
	file, err := ioutil.ReadFile(cfgFilename)
	if err != nil {
		glog.Errorf("Could not read config file %s: %s", cfgFilename, err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		glog.Errorf("Could not marshal config file %s: %s", cfgFilename, err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}

	return config, nil
}
