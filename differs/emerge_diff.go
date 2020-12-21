/*
Copyright 2020 Google, Inc. All rights reserved.

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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
	"github.com/sirupsen/logrus"
)

//Emerge package database location
const emergePkgFile string = "/var/db/pkg"

type EmergeAnalyzer struct{}

func (em EmergeAnalyzer) Name() string {
	return "EmergeAnalyzer"
}

// Diff compares the packages installed by emerge.
func (em EmergeAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff, err := singleVersionDiff(image1, image2, em)
	return diff, err
}

func (em EmergeAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	analysis, err := singleVersionAnalysis(image, em)
	return analysis, err
}

func (em EmergeAnalyzer) getPackages(image pkgutil.Image) (map[string]util.PackageInfo, error) {
	var path string
	if image.FSPath == "" {
		path = emergePkgFile
	} else {
		path = filepath.Join(image.FSPath, emergePkgFile)
	}

	packages := make(map[string]util.PackageInfo)
	if _, err := os.Stat(path); err != nil {
		// invalid image directory path
		logrus.Errorf("Invalid image directory path %s", path)
		return packages, err
	}

	contents, err := ioutil.ReadDir(path)
	if err != nil {
		logrus.Errorf("Non-content in image directory path %s", path)
		return packages, err
	}

	// for i := 0; i < len(contents); i++ {
	for _, c := range contents {
		// c := contents[i]
		pkgPrefix := c.Name()
		pkgContents, err := ioutil.ReadDir(filepath.Join(path, pkgPrefix))
		if err != nil {
			return packages, err
		}
		// for j := 0; j < len(pkgContents); j++ {
		for _, c := range pkgContents {
			// c := pkgContents[j]
			pkgRawName := c.Name()
			// usually, the name of a package installed by emerge is formatted as '{pkgName}-{version}' e.g.(pymongo-3.9.0)
			s := strings.Split(pkgRawName, "-")
			if len(s) != 2 {
				continue
			}
			pkgName, version := s[0], s[1]
			pkgPath := filepath.Join(path, pkgPrefix, pkgRawName, "SIZE")
			size, err := getPkgSize(pkgPath)
			if err != nil {
				return packages, err
			}
			currPackage := util.PackageInfo{Version: version, Size: size}
			fullPackageName := strings.Join([]string{pkgPrefix, pkgName}, "/")
			packages[fullPackageName] = currPackage
		}
	}

	return packages, nil
}

// emerge will count the total size of a package and store it as a SIZE file in pkg metadata directory
// getPkgSize read this SIZE file of a given package
func getPkgSize(pkgPath string) (int64, error) {
	sizeFile, err := os.Open(pkgPath)
	if err != nil {
		logrus.Warnf("unable to open SIZE file for pkg %s", pkgPath)
		return 0, err
	}
	defer sizeFile.Close()
	fileBody, err := ioutil.ReadAll(sizeFile)
	if err != nil {
		logrus.Warnf("unable to read SIZE file for pkg %s", pkgPath)
		return 0, err
	}
	strFileBody := strings.Replace(string(fileBody), "\n", "", -1)
	size, err := strconv.ParseInt(strFileBody, 10, 64)
	if err != nil {
		logrus.Warnf("unable to compute size for pkg %s", pkgPath)
		return 0, err
	}
	return size, nil
}
