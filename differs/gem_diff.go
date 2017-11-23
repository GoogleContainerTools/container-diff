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
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
	"github.com/sirupsen/logrus"
)

type GemAnalyzer struct {
}

func (a GemAnalyzer) Name() string {
	return "GemAnalyzer"
}

// GemAnalyzer.Diff compares gem/bundler installed Ruby packages between layers of two different images
func (a GemAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff, err := multiVersionDiff(image1, image2, a)
	return diff, err
}

func (a GemAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	analysis, err := multiVersionAnalysis(image, a)
	return analysis, err
}

func (a GemAnalyzer) getPackages(image pkgutil.Image) (map[string]map[string]util.PackageInfo, error) {
	path := image.FSPath
	packages := make(map[string]map[string]util.PackageInfo)
	rubyPaths := []string{}
	config, err := image.Image.ConfigFile()
	if err != nil {
		return packages, err
	}
	if config.Config.Env != nil {
		paths := getGemPaths(config.Config.Env)
		for _, p := range paths {
			rubyPaths = append(rubyPaths, p)
		}
	}
	rubyVersions, err := getRubyVersion(path)
	if err != nil {
		logrus.Warnf("Image doesn't appear to have Ruby installed")
		return packages, nil
	}

	// common UNIX ruby installation paths
	for _, rubyVersion := range rubyVersions {
		rubyPaths = append(rubyPaths, filepath.Join(path, "usr/lib/ruby/gems", rubyVersion, "specifications"))
		rubyPaths = append(rubyPaths, filepath.Join(path, "usr/local/lib/ruby/gems", rubyVersion, "specifications"))
	}

	// common "bundler" path
	rubyPaths = append(rubyPaths, filepath.Join(path, "usr/local/bundle/specifications"))

	for _, rubyPath := range rubyPaths {
		contents, err := ioutil.ReadDir(rubyPath)
		if err != nil {
			logrus.Warnf("Could not find Rubygem specifications in %s", rubyPath)
			continue
		}

		for i := 0; i < len(contents); i++ {
			c := contents[i]
			fileName := c.Name()
			// check if gemspec
			packageDir := regexp.MustCompile("^([a-z|A-Z|0-9|_|-]+)-([0-9|\\.|a-z]+)\\.(gemspec)$")
			packageMatch := packageDir.FindStringSubmatch(fileName)
			if len(packageMatch) != 0 {
				packageName := packageMatch[1]
				version := packageMatch[2][:len(packageMatch[2])]

				var size int64
				packagePath := filepath.Join(rubyPath, packageName+"-"+version+".gemspec")
				size = pkgutil.GetSize(packagePath)
				currPackage := util.PackageInfo{Version: version, Size: size}
				mapPath := strings.Replace(rubyPath, path, "", 1)
				addToMap(packages, packageName, mapPath, currPackage)
			}
		}
	}

	return packages, nil
}

func getRubyVersion(pathToLayer string) ([]string, error) {
	matches := []string{}
	pattern := regexp.MustCompile("^[0-9]+\\.[0-9]+\\.[0-9]+$")

	libPaths := []string{"usr/local/lib/ruby", "usr/lib/ruby"}
	for _, lp := range libPaths {
		libPath := filepath.Join(pathToLayer, lp)
		libContents, err := ioutil.ReadDir(libPath)
		if err != nil {
			logrus.Warnf("Could not find %s to determine Ruby version", err)
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

func getGemPaths(vars []string) []string {
	paths := []string{}

	for _, envVar := range vars {
		rubyPathPattern := regexp.MustCompile("^GEM_HOME=(.*)")
		match := rubyPathPattern.FindStringSubmatch(envVar)
		if len(match) != 0 {
			rubyPath := match[1]
			paths = strings.Split(rubyPath, ":")
			break
		}
	}

	return paths
}
