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

package differs

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
	"github.com/sirupsen/logrus"
)

// APK package database location
const apkInstalledDbFile string = "lib/apk/db/installed"

type ApkAnalyzer struct {
}

func (a ApkAnalyzer) Name() string {
	return "ApkAnalyzer"
}

// ApkDiff compares the packages installed by apt-get.
func (a ApkAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff, err := singleVersionDiff(image1, image2, a)
	return diff, err
}

func (a ApkAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	analysis, err := singleVersionAnalysis(image, a)
	return analysis, err
}

func (a ApkAnalyzer) getPackages(image pkgutil.Image) (map[string]util.PackageInfo, error) {
	return apkReadInstalledDbFile(image.FSPath)
}

type ApkLayerAnalyzer struct {
}

func (a ApkLayerAnalyzer) Name() string {
	return "ApkLayerAnalyzer"
}

// ApkDiff compares the packages installed by apt-get.
func (a ApkLayerAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff, err := singleVersionLayerDiff(image1, image2, a)
	return diff, err
}

func (a ApkLayerAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	analysis, err := singleVersionLayerAnalysis(image, a)
	return analysis, err
}

func (a ApkLayerAnalyzer) getPackages(image pkgutil.Image) ([]map[string]util.PackageInfo, error) {
	var packages []map[string]util.PackageInfo
	if _, err := os.Stat(image.FSPath); err != nil {
		// invalid image directory path
		return packages, err
	}
	installedDbFile := filepath.Join(image.FSPath, apkInstalledDbFile)
	if _, err := os.Stat(installedDbFile); err != nil {
		// installed DB file does not exist in this image
		return packages, nil
	}
	for _, layer := range image.Layers {
		layerPackages, err := apkReadInstalledDbFile(layer.FSPath)
		if err != nil {
			return packages, err
		}
		packages = append(packages, layerPackages)
	}

	return packages, nil
}

func apkReadInstalledDbFile(root string) (map[string]util.PackageInfo, error) {
	packages := make(map[string]util.PackageInfo)
	if _, err := os.Stat(root); err != nil {
		// invalid image directory path
		return packages, err
	}
	installedDbFile := filepath.Join(root, apkInstalledDbFile)
	if _, err := os.Stat(installedDbFile); err != nil {
		// installed DB file does not exist in this layer
		return packages, nil
	}
	if file, err := os.Open(installedDbFile); err == nil {
		// make sure it gets closed
		defer file.Close()

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		var currPackage string
		for scanner.Scan() {
			currPackage = apkParseLine(scanner.Text(), currPackage, packages)
		}
	} else {
		return packages, err
	}

	return packages, nil
}

func apkParseLine(text string, currPackage string, packages map[string]util.PackageInfo) string {
	line := strings.Split(text, ":")
	if len(line) == 2 {
		key := line[0]
		value := line[1]

		switch key {
		case "P":
			return value
		case "V":
			if packages[currPackage].Version != "" {
				logrus.Warningln("Multiple versions of same package detected.  Diffing such multi-versioning not yet supported.")
				return currPackage
			}
			currPackageInfo, ok := packages[currPackage]
			if !ok {
				currPackageInfo = util.PackageInfo{}
			}
			currPackageInfo.Version = value
			packages[currPackage] = currPackageInfo
			return currPackage

		case "I":
			currPackageInfo, ok := packages[currPackage]
			if !ok {
				currPackageInfo = util.PackageInfo{}
			}
			var size int64
			var err error
			size, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				logrus.Errorf("Could not get size for %s: %s", currPackage, err)
				size = -1
			}
			// I is in bytes, so *no* conversion needed to keep consistent with the tool's size units
			currPackageInfo.Size = size
			packages[currPackage] = currPackageInfo
			return currPackage
		default:
			return currPackage
		}
	}
	return currPackage
}
