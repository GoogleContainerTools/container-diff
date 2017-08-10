package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/golang/glog"
)

// MultiVersionPackageDiff stores the difference information between two images which could have multi-version packages.
type MultiVersionPackageDiff struct {
	Image1    string
	Packages1 map[string]map[string]PackageInfo
	Image2    string
	Packages2 map[string]map[string]PackageInfo
	InfoDiff  []MultiVersionInfo
}

// MultiVersionInfo stores the information for one multi-version package in two different images.
type MultiVersionInfo struct {
	Package string
	Info1   []PackageInfo
	Info2   []PackageInfo
}

// PackageDiff stores the difference information between two images.
type PackageDiff struct {
	Image1    string
	Packages1 map[string]PackageInfo
	Image2    string
	Packages2 map[string]PackageInfo
	InfoDiff  []Info
}

// Info stores the information for one package in two different images.
type Info struct {
	Package string
	Info1   PackageInfo
	Info2   PackageInfo
}

// PackageInfo stores the specific metadata about a package.
type PackageInfo struct {
	Version string
	Size    string
}

func contains(info1 []PackageInfo, keys1 []string, key string, value PackageInfo) (int, bool) {
	if len(info1) != len(keys1) {
		return 0, false
	}
	for i, currVal := range info1 {
		if !reflect.DeepEqual(currVal, value) {
			continue
		}
		// Check if both global or local installations by trimming img-id and layer hash
		tempPath1 := strings.SplitN(key, "/", 3)
		tempPath2 := strings.SplitN(keys1[i], "/", 3)
		if len(tempPath1) != 3 || len(tempPath2) != 3 {
			continue
		}
		if tempPath1[2] == tempPath2[2] {
			return i, true
		}
	}
	return 0, false
}

func multiVersionDiff(infoDiff []MultiVersionInfo, key string, map1, map2 map[string]PackageInfo) []MultiVersionInfo {
	diff1Possible := []PackageInfo{}
	diff1PossibleKeys := []string{}

	for key1, value1 := range map1 {
		_, ok := map2[key1]
		if !ok {
			diff1Possible = append(diff1Possible, value1)
			diff1PossibleKeys = append(diff1PossibleKeys, key1)
		} else {
			// if key in both maps, means layer hash is the same therefore packages are the same
			delete(map2, key1)
		}
	}

	diff1 := []PackageInfo{}
	diff2 := []PackageInfo{}
	for key2, value2 := range map2 {
		index, ok := contains(diff1Possible, diff1PossibleKeys, key2, value2)
		if !ok {
			diff2 = append(diff2, value2)
		} else {
			if index == 0 {
				diff1Possible = diff1Possible[1:]
				diff1PossibleKeys = diff1PossibleKeys[1:]
			}
			diff1Possible = append(diff1Possible[:index], diff1Possible[index:]...)
			diff1PossibleKeys = append(diff1PossibleKeys[:index], diff1PossibleKeys[index:]...)
		}
	}

	for _, val := range diff1Possible {
		diff1 = append(diff1, val)
	}
	if len(diff1) > 0 || len(diff2) > 0 {
		infoDiff = append(infoDiff, MultiVersionInfo{key, diff1, diff2})
	}
	return infoDiff
}

func checkPackageMapType(map1, map2 interface{}) (reflect.Type, bool, error) {
	// check types and determine multi-version package maps or not
	map1Kind := reflect.ValueOf(map1)
	map2Kind := reflect.ValueOf(map2)
	if map1Kind.Kind() != reflect.Map && map2Kind.Kind() != reflect.Map {
		return nil, false, fmt.Errorf("Package maps were not maps.  Instead were: %s and %s", map1Kind.Kind(), map2Kind.Kind())
	}
	mapType := reflect.TypeOf(map1)
	if mapType != reflect.TypeOf(map2) {
		return nil, false, fmt.Errorf("Package maps were of different types")
	}
	multiVersion := false
	if mapType == reflect.TypeOf(map[string]map[string]PackageInfo{}) {
		multiVersion = true
	}
	return mapType, multiVersion, nil
}

// GetMapDiff determines the differences between maps of package names to PackageInfo structs
// This getter supports only single version packages.
func GetMapDiff(map1, map2 map[string]PackageInfo, img1, img2 string) PackageDiffResult {
	diff := diffMaps(map1, map2)
	diffVal := reflect.ValueOf(diff)
	packDiff := diffVal.Interface().(PackageDiff)
	packDiff.Image1 = img1
	packDiff.Image2 = img2
	return PackageDiffResult{Diff: packDiff}
}

// GetMultiVersionMapDiff determines the differences between two image package maps with multi-version packages
// This getter supports multi version packages.
func GetMultiVersionMapDiff(map1, map2 map[string]map[string]PackageInfo, img1, img2 string) MultiVersionPackageDiffResult {
	diff := diffMaps(map1, map2)
	diffVal := reflect.ValueOf(diff)
	packDiff := diffVal.Interface().(MultiVersionPackageDiff)
	packDiff.Image1 = img1
	packDiff.Image2 = img2
	return MultiVersionPackageDiffResult{Diff: packDiff}
}

// DiffMaps determines the differences between maps of package names to PackageInfo structs
// The return struct includes a list of packages only in the first map, a list of packages only in
// the second map, and a list of packages which differed only in their PackageInfo (version, size, etc.)
func diffMaps(map1, map2 interface{}) interface{} {
	mapType, multiV, err := checkPackageMapType(map1, map2)
	if err != nil {
		glog.Error(err)
	}

	map1Value := reflect.ValueOf(map1)
	map2Value := reflect.ValueOf(map2)

	diff1 := reflect.MakeMap(mapType)
	diff2 := reflect.MakeMap(mapType)
	infoDiff := []Info{}
	multiInfoDiff := []MultiVersionInfo{}

	for _, key1 := range map1Value.MapKeys() {
		value1 := map1Value.MapIndex(key1)
		value2 := map2Value.MapIndex(key1)
		if !value2.IsValid() {
			diff1.SetMapIndex(key1, value1)
		} else if !reflect.DeepEqual(value2.Interface(), value1.Interface()) {
			if multiV {
				multiInfoDiff = multiVersionDiff(multiInfoDiff, key1.String(),
					value1.Interface().(map[string]PackageInfo), value2.Interface().(map[string]PackageInfo))
			} else {
				infoDiff = append(infoDiff, Info{key1.String(), value1.Interface().(PackageInfo),
					value2.Interface().(PackageInfo)})
			}
			map2Value.SetMapIndex(key1, reflect.Value{})
		} else {
			map2Value.SetMapIndex(key1, reflect.Value{})
		}
	}
	for _, key2 := range map2Value.MapKeys() {
		value2 := map2Value.MapIndex(key2)
		diff2.SetMapIndex(key2, value2)
	}
	if multiV {
		return MultiVersionPackageDiff{Packages1: diff1.Interface().(map[string]map[string]PackageInfo),
			Packages2: diff2.Interface().(map[string]map[string]PackageInfo), InfoDiff: multiInfoDiff}
	}
	return PackageDiff{Packages1: diff1.Interface().(map[string]PackageInfo),
		Packages2: diff2.Interface().(map[string]PackageInfo), InfoDiff: infoDiff}
}

func (pi PackageInfo) string() string {
	return pi.Version
}

// BuildLayerTargets creates a string slice of the layers found at path with the target concatenated.
func BuildLayerTargets(path, target string) ([]string, error) {
	layerStems := []string{}
	layers, err := ioutil.ReadDir(path)
	if err != nil {
		return layerStems, err
	}
	for _, layer := range layers {
		if layer.IsDir() {
			layerStems = append(layerStems, filepath.Join(path, layer.Name(), target))
		}
	}
	return layerStems, nil
}
