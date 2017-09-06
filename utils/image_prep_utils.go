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

package utils

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/containers/image/docker"
	"github.com/containers/image/docker/tarfile"
	"github.com/containers/image/types"
	"github.com/docker/docker/client"

	"github.com/golang/glog"
)

var sourceToPrepMap = map[string]func(ip ImagePrepper) Prepper{
	"ID":  func(ip ImagePrepper) Prepper { return IDPrepper{ImagePrepper: ip} },
	"URL": func(ip ImagePrepper) Prepper { return CloudPrepper{ImagePrepper: ip} },
	"tar": func(ip ImagePrepper) Prepper { return TarPrepper{ImagePrepper: ip} },
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

// CloudPrepper prepares images sourced from a Cloud registry
type CloudPrepper struct {
	ImagePrepper
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

	// By default, the image library will try to look at /etc/docker/certs.d
	// As a non-root user, this would result in a permissions error, so we avoid this
	// by looking in a temporary directory we create in the container-diff home directory
	cwd, _ := os.Getwd()
	tmpCerts, _ := ioutil.TempDir(cwd, "certs")
	defer os.RemoveAll(tmpCerts)
	ctx := &types.SystemContext{
		DockerCertPath: tmpCerts,
	}
	img, err := ref.NewImage(ctx)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	defer img.Close()

	imgSrc, err := ref.NewImageSource(ctx, nil)
	if err != nil {
		glog.Error(err)
		return "", err
	}

	for _, b := range img.LayerInfos() {
		bi, _, err := imgSrc.GetBlob(b)
		if err != nil {
			glog.Errorf("Diff may be inaccurate, failed to pull image layer with error: %s", err)
		}
		gzf, err := gzip.NewReader(bi)
		if err != nil {
			glog.Errorf("Diff may be inaccurate, failed to read layers with error: %s", err)
		}
		tr := tar.NewReader(gzf)
		err = unpackTar(tr, path)
		if err != nil {
			glog.Errorf("Diff may be inaccurate, failed to untar layer with error: %s", err)
		}
	}
	return path, nil
}

func (p CloudPrepper) getConfig() (ConfigSchema, error) {
	ref, err := docker.ParseReference("//" + p.Source)
	if err != nil {
		return ConfigSchema{}, err
	}

	// By default, the image library will try to look at /etc/docker/certs.d
	// As a non-root user, this would result in a permissions error, so we avoid this
	// by looking in a temporary directory we create in the container-diff home directory
	cwd, _ := os.Getwd()
	tmpCerts, _ := ioutil.TempDir(cwd, "certs")
	defer os.RemoveAll(tmpCerts)
	ctx := &types.SystemContext{
		DockerCertPath: tmpCerts,
	}
	img, err := ref.NewImage(ctx)
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

type IDPrepper struct {
	ImagePrepper
}

func (p IDPrepper) getFileSystem() (string, error) {
	tarPath, err := saveImageToTar(p.Client, p.Source, p.Source)
	if err != nil {
		return "", err
	}

	defer os.Remove(tarPath)
	return getImageFromTar(tarPath)
}

func (p IDPrepper) getConfig() (ConfigSchema, error) {
	inspect, _, err := p.Client.ImageInspectWithRaw(context.Background(), p.Source)
	if err != nil {
		return ConfigSchema{}, err
	}

	config := ConfigObject{
		Env: inspect.Config.Env,
	}
	history := p.getHistory()
	return ConfigSchema{
		Config:  config,
		History: history,
	}, nil
}

func (p IDPrepper) getHistory() []ImageHistoryItem {
	history, err := p.Client.ImageHistory(context.Background(), p.Source)
	if err != nil {
		glog.Error("Could not obtain image history for %s: %s", p.Source, err)
	}
	historyItems := []ImageHistoryItem{}
	for _, item := range history {
		historyItems = append(historyItems, ImageHistoryItem{CreatedBy: item.CreatedBy})
	}
	return historyItems
}

type TarPrepper struct {
	ImagePrepper
}

func (p TarPrepper) getFileSystem() (string, error) {
	return getImageFromTar(p.Source)
}

func (p TarPrepper) getConfig() (ConfigSchema, error) {
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
