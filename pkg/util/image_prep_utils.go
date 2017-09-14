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
	"archive/tar"
	"compress/gzip"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/container-diff/pkg/util/preppers"
	"github.com/containers/image/types"
	"github.com/docker/docker/client"

	"github.com/golang/glog"
)

var sourceToPrepMap = map[string]func(ip ImagePrepper) Prepper{
	"ID":  func(ip ImagePrepper) Prepper { return preppers.IDPrepper{ImagePrepper: ip} },
	"URL": func(ip ImagePrepper) Prepper { return preppers.CloudPrepper{ImagePrepper: ip} },
	"tar": func(ip ImagePrepper) Prepper { return preppers.TarPrepper{ImagePrepper: ip} },
}

var sourceCheckMap = map[string]func(string) bool{
	"ID":  CheckImageID,
	"URL": CheckImageURL,
	"tar": CheckTar,
}

type Image struct {
	Source string
	FSPath string
	Config ConfigSchema
}

type ImageHistoryItem struct {
	CreatedBy string `json:"created_by"`
}

type ConfigObject struct {
	Env []string `json:"Env"`
}

type ConfigSchema struct {
	Config  ConfigObject       `json:"config"`
	History []ImageHistoryItem `json:"history"`
}

type ImagePrepper struct {
	Source string
	Client *client.Client
}

type Prepper interface {
	getFileSystem() (string, error)
	getConfig() (ConfigSchema, error)
}

func (p ImagePrepper) GetImage() (Image, error) {
	glog.Infof("Starting prep for image %s", p.Source)
	img := p.Source

	var prepper Prepper
	for source, check := range sourceCheckMap {
		if check(img) {
			prepper = sourceToPrepMap[source](p)
			break
		}
	}
	if prepper == nil {
		return Image{}, errors.New("Could not retrieve image from source")
	}

	imgPath, err := prepper.getFileSystem()
	if err != nil {
		return Image{}, err
	}

	config, err := prepper.getConfig()
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

func getImageFromTar(tarPath string) (string, error) {
	glog.Info("Extracting image tar to obtain image file system")
	path := strings.TrimSuffix(tarPath, filepath.Ext(tarPath))
	err := unpackDockerSave(tarPath, path)
	return path, err
}

func CleanupImage(image Image) {
	if image.FSPath != "" {
		glog.Infof("Removing image filesystem directory %s from system", image.FSPath)
		errMsg := remove(image.FSPath, true)
		if errMsg != "" {
			glog.Error(errMsg)
		}
	}
}

func remove(path string, dir bool) string {
	var errStr string
	if path == "" {
		return ""
	}

	var err error
	if dir {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		errStr = "\nUnable to remove " + path
	}
	return errStr
}

func processImageReference(ref types.ImageReference, path string) (string, error) {
	img, err := ref.NewImage(nil)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	defer img.Close()

	imgSrc, err := ref.NewImageSource(nil, nil)
	if err != nil {
		glog.Error(err)
		return "", err
	}

	for _, b := range img.LayerInfos() {
		bi, _, err := imgSrc.GetBlob(b)
		if err != nil {
			glog.Errorf("Failed to pull image layer with error: %s", err)
		}
		gzf, err := gzip.NewReader(bi)
		if err != nil {
			glog.Errorf("Failed to read layers with error: %s", err)
		}
		tr := tar.NewReader(gzf)
		err = unpackTar(tr, path)
		if err != nil {
			glog.Errorf("Failed to untar layer with error: %s", err)
		}
	}
	return path, nil
}
