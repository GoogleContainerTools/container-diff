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
	"strings"

	"github.com/containers/image/docker"
	"github.com/containers/image/docker/reference"
)

const RemotePrefix = "remote://"

// CloudPrepper prepares images sourced from a Cloud registry
type CloudPrepper struct {
	ImagePrepper
}

func (p CloudPrepper) Name() string {
	return "Cloud Registry"
}

func (p CloudPrepper) RawSource() string {
	return strings.Replace(p.Source, RemotePrefix, "", -1)
}

func (p CloudPrepper) SupportsImage() bool {
	// will fail for strings prefixed with 'daemon://'
	if strings.HasPrefix(p.Source, DaemonPrefix) || IsTar(p.Source) {
		return false
	}
	_, err := reference.Parse(p.RawSource())
	return (err == nil)
}

func (p CloudPrepper) GetFileSystem() (string, error) {
	ref, err := docker.ParseReference("//" + p.RawSource())
	if err != nil {
		return "", err
	}

	return getFileSystemFromReference(ref, p.RawSource())
}

func (p CloudPrepper) GetConfig() (ConfigSchema, error) {
	ref, err := docker.ParseReference("//" + p.RawSource())
	if err != nil {
		return ConfigSchema{}, err
	}

	return getConfigFromReference(ref, p.RawSource())
}
