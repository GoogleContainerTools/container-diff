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
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	"github.com/golang/glog"
)

type ImagePrepper struct {
	Source string
	Client *client.Client
}

type Prepper interface {
	getFileSystem() (string, error)
	getConfig() (ConfigSchema, error)
}

func getImage(prepper Prepper, source string) (Image, error) {
	imgPath, err := prepper.getFileSystem()
	if err != nil {
		return Image{}, fmt.Errorf("error msg: %s", err.Error())
	}

	config, err := prepper.getConfig()
	if err != nil {
		return Image{}, fmt.Errorf("error msg: %s", err.Error())
	}

	glog.Infof("Finished prepping image %s", source)
	return Image{
		Source: source,
		FSPath: imgPath,
		Config: config,
	}, nil
}

func (p *ImagePrepper) GetImage() (Image, error) {
	glog.Infof("Starting prep for image %s", p.Source)

	var prepper Prepper

	// first, respect prefixes on image names
	if strings.HasPrefix(p.Source, "daemon://") {
		p.Source = strings.Replace(p.Source, "daemon://", "", -1)
		prepper = DaemonPrepper{ImagePrepper: p}
		return getImage(prepper, p.Source)
	} else if strings.HasPrefix(p.Source, "remote://") {
		p.Source = strings.Replace(p.Source, "remote://", "", -1)
		prepper = CloudPrepper{ImagePrepper: p}
		return getImage(prepper, p.Source)
	}

	// if no prefix found, check local daemon first, otherwise default to cloud registry
	for source, check := range sourceCheckMap {
		if check(p.Source) {
			prepper = sourceToPrepMap[source](p)
			glog.Infof("Attempting to retrieve image with source %s", source)

			img, err := getImage(prepper, p.Source)
			if err != nil {
				glog.Warning(err.Error())
				continue
			}
			return img, nil
		}
	}
	return Image{}, errors.New("Could not retrieve image from source")
}
