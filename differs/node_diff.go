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

package differs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/GoogleCloudPlatform/container-diff/util"
	"github.com/sirupsen/logrus"
)

type NodeAnalyzer struct {
}

func (a NodeAnalyzer) Name() string {
	return "NodeAnalyzer"
}

// NodeDiff compares the packages installed by apt-get.
func (a NodeAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff, err := multiVersionDiff(image1, image2, a)
	return diff, err
}

func (a NodeAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	analysis, err := multiVersionAnalysis(image, a)
	return analysis, err
}

func (a NodeAnalyzer) getPackages(image pkgutil.Image) (map[string]map[string]util.PackageInfo, error) {
	path := image.FSPath
	packages := make(map[string]map[string]util.PackageInfo)
	if _, err := os.Stat(path); err != nil {
		// path provided invalid
		return packages, err
	}
	layerStems, err := buildNodePaths(path)
	if err != nil {
		logrus.Warningf("Error building JSON paths at %s: %s\n", path, err)
		return packages, err
	}

	for _, modulesDir := range layerStems {
		packageJSONs, _ := util.BuildLayerTargets(modulesDir, "package.json")
		for _, currPackage := range packageJSONs {
			if _, err := os.Stat(currPackage); err != nil {
				// package.json file does not exist at this target path
				continue
			}
			packageJSON, err := readPackageJSON(currPackage)
			if err != nil {
				logrus.Warningf("Error reading package JSON at %s: %s\n", currPackage, err)
				return packages, err
			}
			// Build PackageInfo for this package occurence
			var currInfo util.PackageInfo
			currInfo.Version = packageJSON.Version
			packagePath := strings.TrimSuffix(currPackage, "package.json")
			currInfo.Size = pkgutil.GetSize(packagePath)
			mapPath := strings.Replace(packagePath, path, "", 1)
			// Check if other package version already recorded
			if _, ok := packages[packageJSON.Name]; !ok {
				// package not yet seen
				infoMap := make(map[string]util.PackageInfo)
				infoMap[mapPath] = currInfo
				packages[packageJSON.Name] = infoMap
				continue
			}
			packages[packageJSON.Name][mapPath] = currInfo

		}
	}
	return packages, nil
}

type nodePackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func buildNodePaths(path string) ([]string, error) {
	globalPaths := filepath.Join(path, "node_modules")
	localPath := filepath.Join(path, "usr/local/lib/node_modules")
	return []string{globalPaths, localPath}, nil
}

func readPackageJSON(path string) (nodePackage, error) {
	var currPackage nodePackage
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return currPackage, err
	}
	err = json.Unmarshal(jsonBytes, &currPackage)
	if err != nil {
		return currPackage, err
	}
	return currPackage, err
}
