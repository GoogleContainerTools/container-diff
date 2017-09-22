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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/golang/glog"
)

type Event struct {
	Status         string `json:"status"`
	Error          string `json:"error"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func NewClient() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("Error getting docker client: %s", err)
	}
	cli.NegotiateAPIVersion(context.Background())

	return cli, nil
}

func getLayersFromManifest(manifestPath string) ([]string, error) {
	type Manifest struct {
		Layers []string
	}

	manifestJSON, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		errMsg := fmt.Sprintf("Could not open manifest to get layer order: %s", err)
		return []string{}, errors.New(errMsg)
	}

	var imageManifest []Manifest
	err = json.Unmarshal(manifestJSON, &imageManifest)
	if err != nil {
		errMsg := fmt.Sprintf("Could not unmarshal manifest to get layer order: %s", err)
		return []string{}, errors.New(errMsg)
	}
	return imageManifest[0].Layers, nil
}

func unpackDockerSave(tarPath string, target string) error {
	if _, ok := os.Stat(target); ok != nil {
		os.MkdirAll(target, 0777)
	}

	tempLayerDir := target + "-temp"
	err := UnTar(tarPath, tempLayerDir)
	if err != nil {
		errMsg := fmt.Sprintf("Could not unpack saved Docker image %s: %s", tarPath, err)
		return errors.New(errMsg)
	}

	manifest := filepath.Join(tempLayerDir, "manifest.json")
	layers, err := getLayersFromManifest(manifest)
	if err != nil {
		return err
	}

	for _, layer := range layers {
		layerTar := filepath.Join(tempLayerDir, layer)
		if _, err := os.Stat(layerTar); err != nil {
			glog.Infof("Did not unpack layer %s because no layer.tar found", layer)
			continue
		}
		err = UnTar(layerTar, target)
		if err != nil {
			glog.Errorf("Could not unpack layer %s: %s", layer, err)
		}
	}
	err = os.RemoveAll(tempLayerDir)
	if err != nil {
		glog.Errorf("Error deleting temp image layer filesystem: %s", err)
	}
	return nil
}

// ImageToTar writes an image to a .tar file
func saveImageToTar(cli client.APIClient, image, tarName string) (string, error) {
	glog.Info("Saving image")
	imgBytes, err := cli.ImageSave(context.Background(), []string{image})
	if err != nil {
		return "", err
	}
	defer imgBytes.Close()
	newpath := tarName + ".tar"
	return newpath, copyToFile(newpath, imgBytes)
}
