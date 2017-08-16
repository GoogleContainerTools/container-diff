package differs

import (
	"reflect"
	"strings"

	"github.com/GoogleCloudPlatform/container-diff/utils"
)

type MultiVersionPackageAnalyzer interface {
	getPackages(image utils.Image) (map[string]map[string]utils.PackageInfo, error)
}

type SingleVersionPackageAnalyzer interface {
	getPackages(image utils.Image) (map[string]utils.PackageInfo, error)
}

func multiVersionDiff(image1, image2 utils.Image, differ MultiVersionPackageAnalyzer) (utils.DiffResult, error) {
	pack1, err := differ.getPackages(image1)
	if err != nil {
		return &utils.MultiPackageDiffResult{}, err
	}
	pack2, err := differ.getPackages(image2)
	if err != nil {
		return &utils.MultiPackageDiffResult{}, err
	}

	diff := utils.GetMultiVersionMapDiff(pack1, pack2)
	return &utils.MultiPackageDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: strings.TrimSuffix(reflect.TypeOf(differ).Name(), "Analyzer"),
		Diff:     diff,
	}, nil
}

func singleVersionDiff(image1, image2 utils.Image, differ SingleVersionPackageAnalyzer) (utils.DiffResult, error) {
	pack1, err := differ.getPackages(image1)
	if err != nil {
		return &utils.PackageDiffResult{}, err
	}
	pack2, err := differ.getPackages(image2)
	if err != nil {
		return &utils.PackageDiffResult{}, err
	}

	diff := utils.GetMapDiff(pack1, pack2)
	return &utils.PackageDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: strings.TrimSuffix(reflect.TypeOf(differ).Name(), "Analyzer"),
		Diff:     diff,
	}, nil
}

func multiVersionAnalysis(image utils.Image, analyzer MultiVersionPackageAnalyzer) (utils.AnalyzeResult, error) {
	pack, err := analyzer.getPackages(image)
	if err != nil {
		return &utils.MultiPackageAnalyzeResult{}, err
	}

	analysis := utils.MultiPackageAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: strings.TrimSuffix(reflect.TypeOf(analyzer).Name(), "Analyzer"),
		Analysis:    pack,
	}
	return &analysis, nil
}

func singleVersionAnalysis(image utils.Image, analyzer SingleVersionPackageAnalyzer) (utils.AnalyzeResult, error) {
	pack, err := analyzer.getPackages(image)
	if err != nil {
		return &utils.PackageAnalyzeResult{}, err
	}

	analysis := utils.PackageAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: strings.TrimSuffix(reflect.TypeOf(analyzer).Name(), "Analyzer"),
		Analysis:    pack,
	}
	return &analysis, nil
}
