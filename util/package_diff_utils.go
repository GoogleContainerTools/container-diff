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

package util

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"

	"github.com/sirupsen/logrus"
)

// MultiVersionPackageDiff stores the difference information between two images which could have multi-version packages.
type MultiVersionPackageDiff struct {
	Packages1 map[string]map[string]PackageInfo
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
	Packages1 map[string]PackageInfo
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
	Size    int64
}

func multiVersionDiff(infoDiff []MultiVersionInfo, packageName string, map1, map2 map[string]PackageInfo) []MultiVersionInfo {
	diff1 := []PackageInfo{}
	diff2 := []PackageInfo{}
	for path, packInfo1 := range map1 {
		packInfo2, ok := map2[path]
		if !ok {
			diff1 = append(diff1, packInfo1)
			continue
		} else {
			// If a package instance is installed in the same place in Image1 and Image2 with the same version,
			// then they are the same package and should not be included in the diff
			if packInfo1.Version == packInfo2.Version {
				delete(map2, path)
			} else {
				diff1 = append(diff1, packInfo1)
				diff2 = append(diff2, packInfo2)
				delete(map2, path)
			}
		}
	}
	for _, packInfo2 := range map2 {
		diff2 = append(diff2, packInfo2)
	}

	if len(diff1) > 0 || len(diff2) > 0 {
		infoDiff = append(infoDiff, MultiVersionInfo{packageName, diff1, diff2})
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
func GetMapDiff(map1, map2 map[string]PackageInfo) PackageDiff {
	diff := diffMaps(map1, map2)
	diffVal := reflect.ValueOf(diff)
	packDiff := diffVal.Interface().(PackageDiff)
	return packDiff
}

// GetMultiVersionMapDiff determines the differences between two image package maps with multi-version packages
// This getter supports multi version packages.
func GetMultiVersionMapDiff(map1, map2 map[string]map[string]PackageInfo) MultiVersionPackageDiff {
	diff := diffMaps(map1, map2)
	diffVal := reflect.ValueOf(diff)
	packDiff := diffVal.Interface().(MultiVersionPackageDiff)
	return packDiff
}

// DiffMaps determines the differences between maps of package names to PackageInfo structs
// The return struct includes a list of packages only in the first map, a list of packages only in
// the second map, and a list of packages which differed only in their PackageInfo (version, size, etc.)
func diffMaps(map1, map2 interface{}) interface{} {
	mapType, multiV, err := checkPackageMapType(map1, map2)
	if err != nil {
		logrus.Error(err)
	}

	map1Value := reflect.ValueOf(map1)
	map2Value := reflect.ValueOf(map2)

	diff1 := reflect.MakeMap(mapType)
	diff2 := reflect.MakeMap(mapType)
	infoDiff := []Info{}
	multiInfoDiff := []MultiVersionInfo{}

	for _, pack := range map1Value.MapKeys() {
		packageEntry1 := map1Value.MapIndex(pack)
		packageEntry2 := map2Value.MapIndex(pack)
		// If the package does not exist in Image2's map of packages, add it to the collection of unique packages in Image1
		if !packageEntry2.IsValid() {
			diff1.SetMapIndex(pack, packageEntry1)
			// If the package exists in Image2's map of packages but the package information differs between images, add it to
			// the difference.  Otherwise, if the information is consistent across images, delete it from the Image2's map of packages.
		} else {
			if multiV {
				if !reflect.DeepEqual(packageEntry2.Interface(), packageEntry1.Interface()) {
					multiInfoDiff = multiVersionDiff(multiInfoDiff, pack.String(),
						packageEntry1.Interface().(map[string]PackageInfo),
						packageEntry2.Interface().(map[string]PackageInfo))
				}
			} else {
				packageInfo1 := packageEntry1.Interface().(PackageInfo)
				packageInfo2 := packageEntry2.Interface().(PackageInfo)
				// If two instances of the same package don't have the same version, then they are considered to be different
				if packageInfo1.Version != packageInfo2.Version {
					infoDiff = append(infoDiff, Info{pack.String(), packageInfo1, packageInfo2})
				}
			}
			map2Value.SetMapIndex(pack, reflect.Value{})
		}
	}

	// The packages remaining in Image2's map of images are those that exist uniquely in Image2
	for _, key2 := range map2Value.MapKeys() {
		packageEntry2 := map2Value.MapIndex(key2)
		diff2.SetMapIndex(key2, packageEntry2)
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
