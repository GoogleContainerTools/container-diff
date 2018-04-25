/*
Copyright 2018 Google, Inc. All rights reserved.

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
	"github.com/google/go-containerregistry/v1"
	"github.com/google/go-containerregistry/v1/partial"
	"github.com/google/go-containerregistry/v1/types"
)

type TestImage struct {
	Config *v1.ConfigFile
}

func (i *TestImage) RawConfigFile() ([]byte, error) {
	return partial.RawConfigFile(i)
}

func (i *TestImage) ConfigFile() (*v1.ConfigFile, error) {
	return i.Config, nil
}

func (i *TestImage) MediaType() (types.MediaType, error) {
	return types.DockerManifestSchema2, nil
}

func (i *TestImage) LayerByDiffID(diffID v1.Hash) (v1.Layer, error) {
	return nil, nil
}

func (i *TestImage) BlobSet() (map[v1.Hash]struct{}, error) {
	return nil, nil
}

func (i *TestImage) ConfigName() (v1.Hash, error) {
	return v1.Hash{}, nil
}

func (i *TestImage) Digest() (v1.Hash, error) {
	return v1.Hash{}, nil
}

func (i *TestImage) Manifest() (*v1.Manifest, error) {
	return nil, nil
}

func (i *TestImage) RawManifest() ([]byte, error) {
	return nil, nil
}

func (i *TestImage) LayerByDigest(v1.Hash) (v1.Layer, error) {
	return nil, nil
}

func (i *TestImage) Layers() ([]v1.Layer, error) {
	return nil, nil
}
