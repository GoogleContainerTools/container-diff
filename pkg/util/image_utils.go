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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/containers/image/docker/reference"
	"github.com/docker/docker/pkg/system"
	"github.com/golang/glog"
)

func GetImageLayers(pathToImage string) []string {
	layers := []string{}
	contents, err := ioutil.ReadDir(pathToImage)
	if err != nil {
		glog.Error(err.Error())
	}

	for _, file := range contents {
		if file.IsDir() {
			layers = append(layers, file.Name())
		}
	}
	return layers
}

// this method will check all possible legal names for a docker container in a local registry.
// we require that all images passed in have an explicit tag appended to them.
// func CheckImageID(image string) bool {
// 	pattern := regexp.MustCompile("[a-z|0-9]{12}")
// 	if exp := pattern.FindString(image); exp != image {
// 		// didn't find image by ID: let's try [name:tag]
// 		pattern = regexp.MustCompile("[a-z|0-9]*:[a-z|0-9]*")
// 		return pattern.MatchString(image)
// 	}
// 	return true
// }

// func CheckImageURL(image string) bool {
// 	pattern := regexp.MustCompile("^.+/.+(:.+){0,1}$")
// 	if exp := pattern.FindString(image); exp != image || CheckTar(image) {
// 		return false
// 	}
// 	return true
// }

func checkValidImageID(image string) bool {
	_, err := reference.Parse(image)
	return (err == nil)
}

func CheckValidLocalImageID(image string) bool {
	return checkValidImageID(strings.Replace(image, "daemon://", "", -1))
}

func CheckValidRemoteImageID(image string) bool {
	daemonRegex := regexp.MustCompile("daemon://.*")
	if match := daemonRegex.MatchString(image); match {
		return false
	}
	return checkValidImageID(image)
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
