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
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"

	"github.com/docker/docker/pkg/system"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	daemonPrefix = "daemon://"
	remotePrefix = "remote://"

	tagRegexStr = ".*:([^/]+$)"
)

type Layer struct {
	FSPath string
	Digest v1.Hash
}

type Image struct {
	Image  v1.Image
	Source string
	FSPath string
	Digest v1.Hash
	Layers []Layer
}

type ImageHistoryItem struct {
	CreatedBy string `json:"created_by"`
}

// GetImageForName retrieves an image by name alone.
// It does not return layer information, or respect caching.
func GetImageForName(imageName string) (Image, error) {
	return GetImage(imageName, false, "")
}

// GetImage infers the source of an image and retrieves a v1.Image reference to it.
// Once a reference is obtained, it attempts to unpack the v1.Image's reader's contents
// into a temp directory on the local filesystem.
func GetImage(imageName string, includeLayers bool, cacheDir string) (Image, error) {
	logrus.Infof("retrieving image: %s", imageName)
	var img v1.Image
	var err error
	if IsTar(imageName) {
		start := time.Now()
		img, err = tarball.ImageFromPath(imageName, nil)
		if err != nil {
			return Image{}, errors.Wrap(err, "retrieving tar from path")
		}
		elapsed := time.Now().Sub(start)
		logrus.Infof("retrieving image ref from tar took %f seconds", elapsed.Seconds())
	} else if strings.HasPrefix(imageName, daemonPrefix) {
		// remove the daemon prefix
		imageName = strings.Replace(imageName, daemonPrefix, "", -1)

		ref, err := name.ParseReference(imageName, name.WeakValidation)
		if err != nil {
			return Image{}, errors.Wrap(err, "parsing image reference")
		}

		start := time.Now()
		// TODO(nkubala): specify gzip.NoCompression here when functional options are supported
		img, err = daemon.Image(ref, daemon.WithBufferedOpener())
		if err != nil {
			return Image{}, errors.Wrap(err, "retrieving image from daemon")
		}
		elapsed := time.Now().Sub(start)
		logrus.Infof("retrieving local image ref took %f seconds", elapsed.Seconds())
	} else {
		// either has remote prefix or has no prefix, in which case we force remote
		imageName = strings.Replace(imageName, remotePrefix, "", -1)
		ref, err := name.ParseReference(imageName, name.WeakValidation)
		if err != nil {
			return Image{}, errors.Wrap(err, "parsing image reference")
		}
		auth, err := authn.DefaultKeychain.Resolve(ref.Context().Registry)
		if err != nil {
			return Image{}, errors.Wrap(err, "resolving auth")
		}
		start := time.Now()
		img, err = remote.Image(ref, remote.WithAuth(auth), remote.WithTransport(http.DefaultTransport))
		if err != nil {
			return Image{}, errors.Wrap(err, "retrieving remote image")
		}
		elapsed := time.Now().Sub(start)
		logrus.Infof("retrieving remote image ref took %f seconds", elapsed.Seconds())
	}

	// create tempdir and extract fs into it
	var layers []Layer
	if includeLayers {
		start := time.Now()
		imgLayers, err := img.Layers()
		if err != nil {
			return Image{}, errors.Wrap(err, "getting image layers")
		}
		for _, layer := range imgLayers {
			layerStart := time.Now()
			digest, err := layer.Digest()
			path, err := getExtractPathForName(digest.String(), cacheDir)
			if err != nil {
				return Image{
					Layers: layers,
				}, errors.Wrap(err, "getting extract path for layer")
			}
			if err := GetFileSystemForLayer(layer, path, nil); err != nil {
				return Image{
					Layers: layers,
				}, errors.Wrap(err, "getting filesystem for layer")
			}
			layers = append(layers, Layer{
				FSPath: path,
				Digest: digest,
			})
			elapsed := time.Now().Sub(layerStart)
			logrus.Infof("time elapsed retrieving layer: %fs", elapsed.Seconds())
		}
		elapsed := time.Now().Sub(start)
		logrus.Infof("time elapsed retrieving image layers: %fs", elapsed.Seconds())
	}

	imageDigest, err := getImageDigest(img)
	if err != nil {
		return Image{}, err
	}
	path, err := getExtractPathForName(RemoveTag(imageName)+"@"+imageDigest.String(), cacheDir)
	if err != nil {
		return Image{}, err
	}
	// extract fs into provided dir
	if err := GetFileSystemForImage(img, path, nil); err != nil {
		return Image{
			FSPath: path,
			Layers: layers,
		}, errors.Wrap(err, "getting filesystem for image")
	}
	return Image{
		Image:  img,
		Source: imageName,
		FSPath: path,
		Digest: imageDigest,
		Layers: layers,
	}, nil
}

func getExtractPathForName(name string, cacheDir string) (string, error) {
	path := cacheDir
	var err error
	if cacheDir != "" {
		// if cachedir doesn't exist, create it
		if _, err := os.Stat(cacheDir); err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(cacheDir, 0700)
			if err != nil {
				return "", err
			}
			logrus.Infof("caching filesystem at %s", cacheDir)
		}
	} else {
		// otherwise, create tempdir
		logrus.Infof("skipping caching")
		path, err = ioutil.TempDir("", strings.Replace(name, "/", "", -1))
		if err != nil {
			return "", err
		}
	}
	return path, nil
}

func getImageDigest(image v1.Image) (digest v1.Hash, err error) {
	start := time.Now()
	digest, err = image.Digest()
	if err != nil {
		return digest, err
	}
	elapsed := time.Now().Sub(start)
	logrus.Infof("time elapsed retrieving image digest: %fs", elapsed.Seconds())
	return digest, nil
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
	empty, err := DirIsEmpty(root)
	if err != nil {
		return err
	}
	if !empty {
		logrus.Infof("using cached filesystem in %s", root)
		return nil
	}
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
