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
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/docker/docker/pkg/system"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/sirupsen/logrus"
)

const tagRegexStr = ".*:([^/]+$)"

type Layer struct {
	FSPath string
}

type Image struct {
	Image  v1.Image
	Source string
	FSPath string
	Layers []Layer
}

type ImageHistoryItem struct {
	CreatedBy string `json:"created_by"`
}

func CleanupImage(image Image) {
	if image.FSPath != "" {
		logrus.Infof("Removing image filesystem directory %s from system", image.FSPath)
		if err := os.RemoveAll(image.FSPath); err != nil {
			logrus.Warn(err.Error())
		}
	}
	if image.Layers != nil {
		for _, layer := range image.Layers {
			if err := os.RemoveAll(layer.FSPath); err != nil {
				logrus.Warn(err.Error())
			}
		}
	}
}

func SortMap(m map[string]string) string {
	pairs := make([]string, 0)
	for key := range m {
		pairs = append(pairs, fmt.Sprintf("%s:%s", key, m[key]))
	}
	sort.Strings(pairs)
	return strings.Join(pairs, " ")
}

// GetFileSystemForLayer unpacks a layer to local disk
func GetFileSystemForLayer(layer v1.Layer, root string, whitelist []string) error {
	contents, err := layer.Uncompressed()
	if err != nil {
		return err
	}
	return unpackTar(tar.NewReader(contents), root, whitelist)
}

// unpack image filesystem to local disk
// if provided directory is not empty, do nothing
func GetFileSystemForImage(image v1.Image, root string, whitelist []string) error {
	empty, err := DirIsEmpty(root)
	if err != nil {
		return err
	}
	if !empty {
		logrus.Infof("using cached filesystem in %s", root)
		return nil
	}
	if err := unpackTar(tar.NewReader(mutate.Extract(image)), root, whitelist); err != nil {
		return err
	}
	return nil
}

func GetImageLayers(pathToImage string) []string {
	layers := []string{}
	contents, err := ioutil.ReadDir(pathToImage)
	if err != nil {
		logrus.Error(err.Error())
	}

	for _, file := range contents {
		if file.IsDir() {
			layers = append(layers, file.Name())
		}
	}
	return layers
}

// copyToFile writes the content of the reader to the specified file
func copyToFile(outfile string, r io.Reader) error {
	// We use sequential file access here to avoid depleting the standby list
	// on Windows. On Linux, this is a call directly to ioutil.TempFile
	tmpFile, err := system.TempFileSequential(filepath.Dir(outfile), ".docker_temp_")
	if err != nil {
		return err
	}

	tmpPath := tmpFile.Name()

	_, err = io.Copy(tmpFile, r)
	tmpFile.Close()

	if err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err = os.Rename(tmpPath, outfile); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}

// checks to see if an image string contains a tag.
func HasTag(image string) bool {
	tagRegex := regexp.MustCompile(tagRegexStr)
	return tagRegex.MatchString(image)
}

// returns a raw image name with the tag removed
func RemoveTag(image string) string {
	if !HasTag(image) {
		return image
	}
	tagRegex := regexp.MustCompile(tagRegexStr)
	parts := tagRegex.FindStringSubmatch(image)
	tag := parts[len(parts)-1]
	return image[0 : len(image)-len(tag)-1]
}
