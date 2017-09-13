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
	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/GoogleCloudPlatform/container-diff/util"
)

type FileAnalyzer struct {
}

func (a FileAnalyzer) Name() string {
	return "FileAnalyzer"
}

// FileDiff diffs two packages and compares their contents
func (a FileAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff, err := diffImageFiles(image1, image2)
	return &util.DirDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: "File",
		Diff:     diff,
	}, err
}

func (a FileAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	var result util.FileAnalyzeResult

	imgDir, err := pkgutil.GetDirectory(image.FSPath, true)
	if err != nil {
		return result, err
	}

	result.Image = image.Source
	result.AnalyzeType = "File"
	result.Analysis = pkgutil.GetDirectoryEntries(imgDir)
	return &result, err
}

func diffImageFiles(image1, image2 pkgutil.Image) (util.DirDiff, error) {
	img1 := image1.FSPath
	img2 := image2.FSPath

	var diff util.DirDiff

	img1Dir, err := pkgutil.GetDirectory(img1, true)
	if err != nil {
		return diff, err
	}
	img2Dir, err := pkgutil.GetDirectory(img2, true)
	if err != nil {
		return diff, err
	}

	diff, _ = util.DiffDirectory(img1Dir, img2Dir)
	return diff, nil
}
