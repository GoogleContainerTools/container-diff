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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/containers/image/types"

	"github.com/containers/image/pkg/compression"
	"github.com/golang/glog"
)

const LOCAL string = "Local Daemon"
const REMOTE string = "Cloud Registry"
const TAR string = "Tar"

var sourceToPrepMap = map[string]func(ip ImagePrepper) Prepper{
	"Local Daemon":   func(ip ImagePrepper) Prepper { return DaemonPrepper{ImagePrepper: ip} },
	"Cloud Registry": func(ip ImagePrepper) Prepper { return CloudPrepper{ImagePrepper: ip} },
	"Tar":            func(ip ImagePrepper) Prepper { return TarPrepper{ImagePrepper: ip} },
}

// since we iterate through these sequentially, order matters.
// we intentionally check local daemon first. diff.go forces alternate behavior
// by checking image prefix.
var sourceCheckMap = map[string]func(string) bool{
	"Local Daemon":   CheckValidLocalImageID,
	"Cloud Registry": CheckValidRemoteImageID,
	"Tar":            CheckTar,
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

func getFileSystemFromReference(ref types.ImageReference, imageName string) (string, error) {

	// This regex when passed a string creates a list of the form
	// [repourl/image:tag, image:tag, tag] (the tag may or may not be present)
	// URLPattern := regexp.MustCompile("^.+/(.+(:.+){0,1})$")
	// URLMatch := URLPattern.FindStringSubmatch(imageName)
	// Removing the ":" so that the image path name can be <image><tag>
	// sanitizedName := strings.Replace(URLMatch[1], ":", "", -1)
	sanitizedName := strings.Replace(imageName, "/", "", -1)
	sanitizedName = strings.Replace(sanitizedName, ":", "", -1)

	path, err := ioutil.TempDir("", sanitizedName)
	if err != nil {
		return "", err
	}

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
			glog.Errorf("Failed to pull image layer: %s", err)
		}
		// try and detect layer compression
		f, reader, err := compression.DetectCompression(bi)
		if err != nil {
			glog.Errorf("Failed to detect image compression: %s", err)
			return "", err
		}
		if f != nil {
			// decompress if necessary
			reader, err = f(reader)
			if err != nil {
				glog.Errorf("Failed to decompress image: %s", err)
				return "", err
			}
		}
		tr := tar.NewReader(reader)
		err = unpackTar(tr, path)
		if err != nil {
			glog.Errorf("Failed to untar layer with error: %s", err)
		}
	}
	return path, nil
}

func getConfigFromReference(ref types.ImageReference, source string) (ConfigSchema, error) {
	img, err := ref.NewImage(nil)
	if err != nil {
		glog.Errorf("Error referencing image %s from registry: %s", source, err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}
	defer img.Close()

	configBlob, err := img.ConfigBlob()
	if err != nil {
		glog.Errorf("Error obtaining config blob for image %s from registry: %s", source, err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}

	var config ConfigSchema
	err = json.Unmarshal(configBlob, &config)
	if err != nil {
		glog.Errorf("Error with config file struct for image %s: %s", source, err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}
	return config, nil
}
