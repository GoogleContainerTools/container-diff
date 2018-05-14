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
	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
	"io/ioutil"
	"os"
)

type FileAnalyzer struct {
}

func (a FileAnalyzer) Name() string {
	return "FileAnalyzer"
}

// FileDiff diffs two packages and compares their contents
func (a FileAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff, err := diffImageFiles(image1.FSPath, image2.FSPath)
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

func diffImageFiles(img1, img2 string) (util.DirDiff, error) {
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

type FileLayerAnalyzer struct {
}

func (a FileLayerAnalyzer) Name() string {
	return "FileLayerAnalyzer"
}

// FileDiff diffs two packages and compares their contents
func (a FileLayerAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	var dirDiffs []util.DirDiff
	// This path to an empty dir will be used for diffing in cases
	// where one image has more layers than the other
	emptyPath, err := ioutil.TempDir("", "")
	if err != nil {
		return &util.MultipleDirDiffResult{}, err
	}
	defer os.RemoveAll(emptyPath)

	// Go through each layer of the first image...
	for index, layer := range image1.Layers {
		// ...if there is no corresponding layer in the second image, diff with the empty dir
		if index >= len(image2.Layers) {
			diff, err := diffImageFiles(layer.FSPath, emptyPath)
			if err != nil {
				return &util.MultipleDirDiffResult{}, err
			}
			dirDiffs = append(dirDiffs, diff)
			continue
		}
		// ...else, diff as usual
		layer2 := image2.Layers[index]
		diff, err := diffImageFiles(layer.FSPath, layer2.FSPath)
		if err != nil {
			return &util.MultipleDirDiffResult{}, err
		}
		dirDiffs = append(dirDiffs, diff)
	}

	// check if there are any additional layers in image2...
	if len(image2.Layers) > len(image1.Layers) {
		// ... and diff any additional layers with the empty dir
		for index := len(image1.Layers); index < len(image2.Layers); index++ {
			layer2 := image2.Layers[index]
			diff, err := diffImageFiles(emptyPath, layer2.FSPath)
			if err != nil {
				return &util.MultipleDirDiffResult{}, err
			}
			dirDiffs = append(dirDiffs, diff)
		}
	}
	return &util.MultipleDirDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: "FileLayer",
		Diff: util.MultipleDirDiff{
			DirDiffs: dirDiffs,
		},
	}, nil
}

func (a FileLayerAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	var directoryEntries [][]pkgutil.DirectoryEntry
	for _, layer := range image.Layers {
		layerDir, err := pkgutil.GetDirectory(layer.FSPath, true)
		if err != nil {
			return util.FileLayerAnalyzeResult{}, err
		}
		directoryEntry := pkgutil.GetDirectoryEntries(layerDir)
		directoryEntries = append(directoryEntries, directoryEntry)
	}

	return &util.FileLayerAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: "FileLayer",
		Analysis:    directoryEntries,
	}, nil
}
