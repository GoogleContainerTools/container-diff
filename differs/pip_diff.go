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
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
	"github.com/sirupsen/logrus"
)

type PipAnalyzer struct {
}

func (a PipAnalyzer) Name() string {
	return "PipAnalyzer"
}

// PipDiff compares pip-installed Python packages between layers of two different images.
func (a PipAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff, err := multiVersionDiff(image1, image2, a)
	return diff, err
}

func (a PipAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	analysis, err := multiVersionAnalysis(image, a)
	return analysis, err
}

func (a PipAnalyzer) getPackages(image pkgutil.Image) (map[string]map[string]util.PackageInfo, error) {
	path := image.FSPath
	packages := make(map[string]map[string]util.PackageInfo)
	pythonPaths := []string{}
	config, err := image.Image.ConfigFile()
	if err != nil {
		return packages, err
	}
	if config.Config.Env != nil {
		paths := getPythonPaths(config.Config.Env)
		for _, p := range paths {
			pythonPaths = append(pythonPaths, p)
		}
	}
	pythonVersions, err := getPythonVersion(path)
	if err != nil {
		// Image doesn't have Python installed
		return packages, nil
	}

	// default python package installation directories in unix
	// these are hardcoded in the python source; unfortunately no way to retrieve them from env
	for _, pythonVersion := range pythonVersions {
		pythonPaths = append(pythonPaths, filepath.Join(path, "usr/lib", pythonVersion))
		pythonPaths = append(pythonPaths, filepath.Join(path, "usr/lib", pythonVersion, "dist-packages"))
		pythonPaths = append(pythonPaths, filepath.Join(path, "usr/lib", pythonVersion, "site-packages"))
		pythonPaths = append(pythonPaths, filepath.Join(path, "usr/local/lib", pythonVersion, "dist-packages"))
		pythonPaths = append(pythonPaths, filepath.Join(path, "usr/local/lib", pythonVersion, "site-packages"))
	}

	for _, pythonPath := range pythonPaths {
		contents, err := ioutil.ReadDir(pythonPath)
		if err != nil {
			// python version folder doesn't have a site-packages folder
			continue
		}

		for i := 0; i < len(contents); i++ {
			c := contents[i]
			fileName := c.Name()
			// check if package
			packageDir := regexp.MustCompile("^([a-z|A-Z|0-9|_]+)-(([0-9]+?\\.){2,3})(dist-info|egg-info)$")
			packageMatch := packageDir.FindStringSubmatch(fileName)
			if len(packageMatch) != 0 {
				packageName := packageMatch[1]
				version := packageMatch[2][:len(packageMatch[2])-1]

				// Retrieves size for actual package/script corresponding to each dist-info metadata directory
				// by taking the file entry alphabetically before it (for a package) or after it (for a script)
				var size int64
				if i-1 >= 0 && contents[i-1].Name() == packageName {
					packagePath := filepath.Join(pythonPath, packageName)
					size = pkgutil.GetSize(packagePath)
				} else if i+1 < len(contents) && contents[i+1].Name() == packageName+".py" {
					size = contents[i+1].Size()
				} else {
					logrus.Errorf("Could not find Python package %s for corresponding metadata info", packageName)
					continue
				}
				currPackage := util.PackageInfo{Version: version, Size: size}
				mapPath := strings.Replace(pythonPath, path, "", 1)
				addToMap(packages, packageName, mapPath, currPackage)
			}
		}
	}

	return packages, nil
}

func addToMap(packages map[string]map[string]util.PackageInfo, pack string, path string, packInfo util.PackageInfo) {
	if _, ok := packages[pack]; !ok {
		// package not yet seen
		infoMap := make(map[string]util.PackageInfo)
		infoMap[path] = packInfo
		packages[pack] = infoMap
		return
	}
	packages[pack][path] = packInfo
}

func getPythonVersion(pathToLayer string) ([]string, error) {
	matches := []string{}
	pattern := regexp.MustCompile("^python[0-9]+\\.[0-9]+$")

	libPaths := []string{"usr/local/lib", "usr/lib"}
	for _, lp := range libPaths {
		libPath := filepath.Join(pathToLayer, lp)
		libContents, err := ioutil.ReadDir(libPath)
		if err != nil {
			logrus.Debugf("Could not find %s to determine Python version", err)
			continue
		}
		for _, file := range libContents {
			match := pattern.FindString(file.Name())
			if match != "" {
				matches = append(matches, match)
			}
		}
	}
	return matches, nil
}

func getPythonPaths(vars []string) []string {
	paths := []string{}
	for _, envVar := range vars {
		pythonPathPattern := regexp.MustCompile("^PYTHONPATH=(.*)")
		match := pythonPathPattern.FindStringSubmatch(envVar)
		if len(match) != 0 {
			pythonPath := match[1]
			paths = strings.Split(pythonPath, ":")
			break
		}
	}
	return paths
}
