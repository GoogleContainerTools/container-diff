package differs

import (
	"reflect"

	"github.com/GoogleCloudPlatform/container-diff/utils"
)

type MultiVersionPackageDiffer interface {
	getPackages(image utils.Image) (map[string]map[string]utils.PackageInfo, error)
}

type SingleVersionPackageDiffer interface {
	getPackages(image utils.Image) (map[string]utils.PackageInfo, error)
}

func multiVersionDiff(image1, image2 utils.Image, differ MultiVersionPackageDiffer) (utils.DiffResult, error) {
	pack1, err := differ.getPackages(image1)
	if err != nil {
		return &utils.MultiVersionPackageDiffResult{}, err
	}
	pack2, err := differ.getPackages(image2)
	if err != nil {
		return &utils.MultiVersionPackageDiffResult{}, err
	}

	diff := utils.GetMultiVersionMapDiff(pack1, pack2, image1.Source, image2.Source)
	diff.DiffType = reflect.TypeOf(differ).Name()
	return &diff, nil
}

func singleVersionDiff(image1, image2 utils.Image, differ SingleVersionPackageDiffer) (utils.DiffResult, error) {
	pack1, err := differ.getPackages(image1)
	if err != nil {
		return &utils.PackageDiffResult{}, err
	}
	pack2, err := differ.getPackages(image2)
	if err != nil {
		return &utils.PackageDiffResult{}, err
	}

	diff := utils.GetMapDiff(pack1, pack2, image1.Source, image2.Source)
	diff.DiffType = reflect.TypeOf(differ).Name()
	return &diff, nil
}
