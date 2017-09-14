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

package preppers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/containers/image/docker"
	"github.com/golang/glog"
)

// CloudPrepper prepares images sourced from a Cloud registry
type CloudPrepper struct {
	util.ImagePrepper
}

func (p CloudPrepper) getFileSystem() (string, error) {
	// The regexp when passed a string creates a list of the form
	// [repourl/image:tag, image:tag, tag] (the tag may or may not be present)
	URLPattern := regexp.MustCompile("^.+/(.+(:.+){0,1})$")
	URLMatch := URLPattern.FindStringSubmatch(p.Source)
	// Removing the ":" so that the image path name can be <image><tag>
	sanitizedName := strings.Replace(URLMatch[1], ":", "", -1)

	path, err := ioutil.TempDir("", sanitizedName)
	if err != nil {
		return "", err
	}

	ref, err := docker.ParseReference("//" + p.Source)
	if err != nil {
		return "", err
	}

	return processImageReference(ref, path)
}

func (p CloudPrepper) getConfig() (ConfigSchema, error) {
	ref, err := docker.ParseReference("//" + p.Source)
	if err != nil {
		return ConfigSchema{}, err
	}

	img, err := ref.NewImage(nil)
	if err != nil {
		glog.Errorf("Error referencing image %s from registry: %s", p.Source, err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}
	defer img.Close()

	configBlob, err := img.ConfigBlob()
	if err != nil {
		glog.Errorf("Error obtaining config blob for image %s from registry: %s", p.Source, err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}

	var config ConfigSchema
	err = json.Unmarshal(configBlob, &config)
	if err != nil {
		glog.Errorf("Error with config file struct for image %s: %s", p.Source, err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}
	return config, nil
}
