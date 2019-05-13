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
	"io/ioutil"
	"os"
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
			var metadata *os.File
			var err error
			if strings.HasSuffix(fileName, "egg-info") {
				// wheel directory
				metadata, err = os.Open(filepath.Join(pythonPath, fileName, "PKG-INFO"))
				if err != nil {
					logrus.Debugf("unable to open PKG-INFO for egg %s", fileName)
				}
			} else if strings.HasSuffix(fileName, "dist-info") {
				// egg directory
				metadata, err = os.Open(filepath.Join(pythonPath, fileName, "METADATA"))
				if err != nil {
					logrus.Debugf("unable to open METADATA for wheel %s", fileName)
				}
			} else {
				// no match
				continue
			}

			var line, packageName, version string
			if metadata == nil {
				// unable to open metadata file: try reading the package itself
				mPath := filepath.Join(pythonPath, fileName)
				metadata, err = os.Open(mPath)
				fInfo, _ := os.Stat(mPath)
				if err != nil || fInfo.IsDir() {
					// if this also doesn't work, the package doesn't have the correct metadata structure
					// try and parse the name using a regex anyway
					logrus.Debugf("failed to locate package metadata: attempting to infer package name")
					packageDir := regexp.MustCompile("^([a-z|A-Z|0-9|_]+)-(([0-9]+?\\.){2,3})(dist-info|egg-info)$")
					packageMatch := packageDir.FindStringSubmatch(fileName)
					if len(packageMatch) != 0 {
						packageName = packageMatch[1]
						version = packageMatch[2][:len(packageMatch[2])-1]
					}
				}
			}

			if metadata != nil {
				scanner := bufio.NewScanner(metadata)
				scanner.Split(bufio.ScanLines)
				for scanner.Scan() {
					line = scanner.Text()
					if strings.HasPrefix(line, "Name") {
						packageName = strings.Split(line, ": ")[1]
						// next line is always the version
						scanner.Scan()
						version = strings.Split(scanner.Text(), ": ")[1]
						break
					}
				}
			}

			// First, try and use the "top_level.txt",
			// Many egg packages contains a "top_level.txt" file describing the directories containing the
			// required code. Combining the sizes of each of these directories should give the total size.
			var size int64
			topLevelReader, err := os.Open(filepath.Join(pythonPath, fileName, "top_level.txt"))
			if err == nil {
				scanner := bufio.NewScanner(topLevelReader)
				scanner.Split(bufio.ScanLines)
				for scanner.Scan() {
					// check if directory exists first, then retrieve size
					contentPath := filepath.Join(pythonPath, scanner.Text())
					if _, err := os.Stat(contentPath); err == nil {
						size = size + pkgutil.GetSize(contentPath)
					} else if _, err := os.Stat(contentPath + ".py"); err == nil {
						// sometimes the top level content is just a single python file; try this too
						size = size + pkgutil.GetSize(contentPath+".py")
					}
				}
			} else {
				logrus.Debugf("unable to use top_level.txt: falling back to alphabetical directory entry heuristic...")

				// Retrieves size for actual package/script corresponding to each dist-info metadata directory
				// by examining the file entries directly before and after it
				if i-1 >= 0 && strings.Contains(contents[i-1].Name(), packageName) {
					packagePath := filepath.Join(pythonPath, contents[i-1].Name())
					size = pkgutil.GetSize(packagePath)
				} else if i+1 < len(contents) && strings.Contains(contents[i+1].Name(), packageName) {
					packagePath := filepath.Join(pythonPath, contents[i+1].Name())
					size = pkgutil.GetSize(packagePath)
				} else {
					logrus.Errorf("failed to locate python package for corresponding package metadata %s", packageName)
					continue
				}
			}

			currPackage := util.PackageInfo{Version: version, Size: size}
			mapPath := strings.Replace(pythonPath, path, "", 1)
			addToMap(packages, packageName, mapPath, currPackage)
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
