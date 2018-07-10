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
	"strings"

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
	"github.com/sirupsen/logrus"
)

type MultiVersionPackageAnalyzer interface {
	getPackages(image pkgutil.Image) (map[string]map[string]util.PackageInfo, error)
	Name() string
}

type SingleVersionPackageAnalyzer interface {
	getPackages(image pkgutil.Image) (map[string]util.PackageInfo, error)
	Name() string
}

type SingleVersionPackageLayerAnalyzer interface {
	getPackages(image pkgutil.Image) ([]map[string]util.PackageInfo, error)
	Name() string
}

func multiVersionDiff(image1, image2 pkgutil.Image, differ MultiVersionPackageAnalyzer) (*util.MultiVersionPackageDiffResult, error) {
	pack1, err := differ.getPackages(image1)
	if err != nil {
		return &util.MultiVersionPackageDiffResult{}, err
	}
	pack2, err := differ.getPackages(image2)
	if err != nil {
		return &util.MultiVersionPackageDiffResult{}, err
	}

	diff := util.GetMultiVersionMapDiff(pack1, pack2)
	return &util.MultiVersionPackageDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: strings.TrimSuffix(differ.Name(), "Analyzer"),
		Diff:     diff,
	}, nil
}

func singleVersionDiff(image1, image2 pkgutil.Image, differ SingleVersionPackageAnalyzer) (*util.SingleVersionPackageDiffResult, error) {
	pack1, err := differ.getPackages(image1)
	if err != nil {
		return &util.SingleVersionPackageDiffResult{}, err
	}
	pack2, err := differ.getPackages(image2)
	if err != nil {
		return &util.SingleVersionPackageDiffResult{}, err
	}

	diff := util.GetMapDiff(pack1, pack2)
	return &util.SingleVersionPackageDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: strings.TrimSuffix(differ.Name(), "Analyzer"),
		Diff:     diff,
	}, nil
}

// singleVersionLayerDiff diffs the packages of each image layer by layer
func singleVersionLayerDiff(image1, image2 pkgutil.Image, differ SingleVersionPackageLayerAnalyzer) (*util.SingleVersionPackageLayerDiffResult, error) {
	pack1, err := differ.getPackages(image1)
	if err != nil {
		return &util.SingleVersionPackageLayerDiffResult{}, err
	}
	pack2, err := differ.getPackages(image2)
	if err != nil {
		return &util.SingleVersionPackageLayerDiffResult{}, err
	}
	var pkgDiffs []util.PackageDiff

	// Go through each layer for image1
	for i := range pack1 {
		if i >= len(pack2) {
			// Skip diff when there is no layer to compare with in image2
			continue
		}

		pkgDiff := util.GetMapDiff(pack1[i], pack2[i])
		pkgDiffs = append(pkgDiffs, pkgDiff)
	}

	if len(image1.Layers) != len(image2.Layers) {
		logrus.Infof("%s and %s have different number of layers, please consider using container-diff analyze to view the contents of each image in each layer", image1.Source, image2.Source)
	}

	return &util.SingleVersionPackageLayerDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: strings.TrimSuffix(differ.Name(), "Analyzer"),
		Diff: util.PackageLayerDiff{
			PackageDiffs: pkgDiffs,
		},
	}, nil
}

func multiVersionAnalysis(image pkgutil.Image, analyzer MultiVersionPackageAnalyzer) (*util.MultiVersionPackageAnalyzeResult, error) {
	pack, err := analyzer.getPackages(image)
	if err != nil {
		return &util.MultiVersionPackageAnalyzeResult{}, err
	}

	analysis := util.MultiVersionPackageAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: strings.TrimSuffix(analyzer.Name(), "Analyzer"),
		Analysis:    pack,
	}
	return &analysis, nil
}

func singleVersionAnalysis(image pkgutil.Image, analyzer SingleVersionPackageAnalyzer) (*util.SingleVersionPackageAnalyzeResult, error) {
	pack, err := analyzer.getPackages(image)
	if err != nil {
		return &util.SingleVersionPackageAnalyzeResult{}, err
	}

	analysis := util.SingleVersionPackageAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: strings.TrimSuffix(analyzer.Name(), "Analyzer"),
		Analysis:    pack,
	}
	return &analysis, nil
}

// singleVersionLayerAnalysis returns the packages included, deleted or
// updated in each layer
func singleVersionLayerAnalysis(image pkgutil.Image, analyzer SingleVersionPackageLayerAnalyzer) (*util.SingleVersionPackageLayerAnalyzeResult, error) {
	pack, err := analyzer.getPackages(image)
	if err != nil {
		return &util.SingleVersionPackageLayerAnalyzeResult{}, err
	}
	var pkgDiffs []util.PackageDiff

	// Each layer with modified packages includes a complete list of packages
	// in its package database. Thus we diff the current layer with the
	// previous one. Not all layers may include differences in packages, those
	// are omitted.
	preInd := -1
	for i := range pack {
		var pkgDiff util.PackageDiff
		if preInd < 0 && len(pack[i]) > 0 {
			pkgDiff = util.GetMapDiff(make(map[string]util.PackageInfo), pack[i])
			preInd = i
		} else if preInd >= 0 && len(pack[i]) > 0 {
			pkgDiff = util.GetMapDiff(pack[preInd], pack[i])
			preInd = i
		}

		pkgDiffs = append(pkgDiffs, pkgDiff)
	}

	return &util.SingleVersionPackageLayerAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: strings.TrimSuffix(analyzer.Name(), "Analyzer"),
		Analysis: util.PackageLayerDiff{
			PackageDiffs: pkgDiffs,
		},
	}, nil
}
