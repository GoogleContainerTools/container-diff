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
	"errors"

	"github.com/GoogleCloudPlatform/container-diff/pkg/cache"
	"github.com/docker/docker/client"
	"github.com/golang/glog"
)

type ImagePrepper struct {
	Source string
	Client *client.Client
	Cache  *cache.FileCache
}

type Prepper interface {
	Name() string
	GetSource() string
	GetFileSystem() (string, error)
	GetConfig() (ConfigSchema, error)
	SupportsImage() bool
}

func (p ImagePrepper) GetImage() (Image, error) {
	glog.Infof("Starting prep for image %s", p.Source)
	img := p.Source

	var prepper Prepper

	for _, prepperConstructor := range orderedPreppers {
		prepper = prepperConstructor(p)
		if prepper.SupportsImage() {
			break
		}
	}

	if prepper == nil {
		return Image{}, errors.New("Could not retrieve image from source")
	}

	imgPath, err := prepper.GetFileSystem()
	if err != nil {
		return Image{}, err
	}

	config, err := prepper.GetConfig()
	if err != nil {
		glog.Error("Error retrieving History: ", err)
	}

	glog.Infof("Finished prepping image %s", p.Source)
	return Image{
		Source: img,
		FSPath: imgPath,
		Config: config,
	}, nil
}
