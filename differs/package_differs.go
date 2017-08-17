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
		return &utils.MultiVersionPackageDiffResult{}, err
	}
	pack2, err := differ.getPackages(image2)
	if err != nil {
		return &utils.MultiVersionPackageDiffResult{}, err
	}

	diff := utils.GetMultiVersionMapDiff(pack1, pack2)
	return &utils.MultiVersionPackageDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: strings.TrimSuffix(reflect.TypeOf(differ).Name(), "Analyzer"),
		Diff:     diff,
	}, nil
}

func singleVersionDiff(image1, image2 utils.Image, differ SingleVersionPackageAnalyzer) (utils.DiffResult, error) {
	pack1, err := differ.getPackages(image1)
	if err != nil {
		return &utils.SingleVersionPackageDiffResult{}, err
	}
	pack2, err := differ.getPackages(image2)
	if err != nil {
		return &utils.SingleVersionPackageDiffResult{}, err
	}

	diff := utils.GetMapDiff(pack1, pack2)
	return &utils.SingleVersionPackageDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: strings.TrimSuffix(reflect.TypeOf(differ).Name(), "Analyzer"),
		Diff:     diff,
	}, nil
}

func multiVersionAnalysis(image utils.Image, analyzer MultiVersionPackageAnalyzer) (utils.AnalyzeResult, error) {
	pack, err := analyzer.getPackages(image)
	if err != nil {
		return &utils.MultiVersionPackageAnalyzeResult{}, err
	}

	analysis := utils.MultiVersionPackageAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: strings.TrimSuffix(reflect.TypeOf(analyzer).Name(), "Analyzer"),
		Analysis:    pack,
	}
	return &analysis, nil
}

func singleVersionAnalysis(image utils.Image, analyzer SingleVersionPackageAnalyzer) (utils.AnalyzeResult, error) {
	pack, err := analyzer.getPackages(image)
	if err != nil {
		return &utils.SingleVersionPackageAnalyzeResult{}, err
	}

	analysis := utils.SingleVersionPackageAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: strings.TrimSuffix(reflect.TypeOf(analyzer).Name(), "Analyzer"),
		Analysis:    pack,
	}
	return &analysis, nil
}
