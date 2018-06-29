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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
)

const gemspecMatchStr = "[a-z|A-Z|0-9|_|-]+-[0-9\\.]+\\.gemspec$"
const gemspecCaptureStr = "^([a-z|A-Z|0-9|_|-]+)-([0-9|\\.|a-z]+)\\.gemspec$"

var gemspecMatch = regexp.MustCompile(gemspecMatchStr)
var gemspecCapture = regexp.MustCompile(gemspecCaptureStr)

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
	basePath := image.FSPath
	packages := make(map[string]map[string]util.PackageInfo)
	// Walk the image's file system
	err := filepath.Walk(basePath, func(searchPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Match if current path is a file with a valid gem name, version, and gemspec extension
		if !f.IsDir() && gemspecMatch.MatchString(searchPath) {
			// use the capture the gem's name and version
			match := gemspecCapture.FindStringSubmatch(f.Name())

			if len(match) != 0 {
				// match[1] => ruby_gem-name
				gemName := match[1]
				// match[2] => 0.1.0
				gemVersion := match[2]

				var size int64
				var path string

				size = pkgutil.GetSize(searchPath)
				currPackage := util.PackageInfo{Version: gemVersion, Size: size}

				// trim path to file system
				path = strings.TrimPrefix(searchPath, basePath)
				// trim package file name
				path = strings.TrimSuffix(path, f.Name())
				addToMap(packages, gemName, path, currPackage)
			}

		}

		return err
	})

	if err != nil {
		return packages, err
	}

	return packages, nil
}
