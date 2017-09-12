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
	"strings"

	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/GoogleCloudPlatform/container-diff/util"
)

type MultiVersionPackageAnalyzer interface {
	getPackages(image pkgutil.Image) (map[string]map[string]util.PackageInfo, error)
	Name() string
}

type SingleVersionPackageAnalyzer interface {
	getPackages(image pkgutil.Image) (map[string]util.PackageInfo, error)
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
