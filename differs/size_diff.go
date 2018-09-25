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
	"strconv"

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
)

type SizeAnalyzer struct {
}

func (a SizeAnalyzer) Name() string {
	return "SizeAnalyzer"
}

// SizeDiff diffs two images and compares their size
func (a SizeAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff := []util.SizeDiff{}
	size1 := pkgutil.GetSize(image1.FSPath)
	size2 := pkgutil.GetSize(image2.FSPath)

	if size1 != size2 {
		diff = append(diff, util.SizeDiff{
			Size1: size1,
			Size2: size2,
		})
	}

	return &util.SizeDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: "Size",
		Diff:     diff,
	}, nil
}

func (a SizeAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	entries := []util.SizeEntry{
		{
			Name:   image.Source,
			Digest: image.Digest,
			Size:   pkgutil.GetSize(image.FSPath),
		},
	}

	return &util.SizeAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: "Size",
		Analysis:    entries,
	}, nil
}

type SizeLayerAnalyzer struct {
}

func (a SizeLayerAnalyzer) Name() string {
	return "SizeLayerAnalyzer"
}

// SizeLayerDiff diffs the layers of two images and compares their size
func (a SizeLayerAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	var layerDiffs []util.SizeDiff

	maxLayer := len(image1.Layers)
	if len(image2.Layers) > maxLayer {
		maxLayer = len(image2.Layers)
	}

	for index := 0; index < maxLayer; index++ {
		var size1, size2 int64 = -1, -1
		if index < len(image1.Layers) {
			size1 = pkgutil.GetSize(image1.Layers[index].FSPath)
		}
		if index < len(image2.Layers) {
			size2 = pkgutil.GetSize(image2.Layers[index].FSPath)
		}

		if size1 != size2 {
			diff := util.SizeDiff{
				Name:  strconv.Itoa(index),
				Size1: size1,
				Size2: size2,
			}
			layerDiffs = append(layerDiffs, diff)
		}
	}

	return &util.SizeLayerDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: "SizeLayer",
		Diff:     layerDiffs,
	}, nil
}

func (a SizeLayerAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	var entries []util.SizeEntry
	for index, layer := range image.Layers {
		entry := util.SizeEntry{
			Name:   strconv.Itoa(index),
			Digest: layer.Digest,
			Size:   pkgutil.GetSize(layer.FSPath),
		}
		entries = append(entries, entry)
	}

	return &util.SizeLayerAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: "SizeLayer",
		Analysis:    entries,
	}, nil
}
