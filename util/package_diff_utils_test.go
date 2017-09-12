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
	"reflect"
	"sort"
	"testing"
)

type ByPackage []Info

func (a ByPackage) Len() int           { return len(a) }
func (a ByPackage) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPackage) Less(i, j int) bool { return a[i].Package < a[j].Package }

type ByMultiPackage []MultiVersionInfo

func (a ByMultiPackage) Len() int           { return len(a) }
func (a ByMultiPackage) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMultiPackage) Less(i, j int) bool { return a[i].Package < a[j].Package }

type ByPackageInfo []PackageInfo

func (a ByPackageInfo) Len() int           { return len(a) }
func (a ByPackageInfo) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPackageInfo) Less(i, j int) bool { return a[i].Version < a[j].Version }

func TestDiffMaps(t *testing.T) {
	testCases := []struct {
		descrip  string
		map1     interface{}
		map2     interface{}
		expected interface{}
	}{
		{
			descrip: "Missing Packages.",
			map1: map[string]PackageInfo{
				"pac1": {"1.0", 40},
				"pac3": {"3.0", 60}},
			map2: map[string]PackageInfo{
				"pac4": {"4.0", 70},
				"pac5": {"5.0", 80}},
			expected: PackageDiff{
				Packages1: map[string]PackageInfo{
					"pac1": {"1.0", 40},
					"pac3": {"3.0", 60}},
				Packages2: map[string]PackageInfo{
					"pac4": {"4.0", 70},
					"pac5": {"5.0", 80}},
				InfoDiff: []Info{}},
		},
		{
			descrip: "Different Versions and Sizes.",
			map1: map[string]PackageInfo{
				"pac2": {"2.0", 50},
				"pac3": {"3.0", 60}},
			map2: map[string]PackageInfo{
				"pac2": {"2.0", 45},
				"pac3": {"4.0", 60}},
			expected: PackageDiff{
				Packages1: map[string]PackageInfo{},
				Packages2: map[string]PackageInfo{},
				InfoDiff: []Info{
					{"pac3", PackageInfo{"3.0", 60}, PackageInfo{"4.0", 60}}},
			},
		},
		{
			descrip: "Identical packages, versions, and sizes",
			map1: map[string]PackageInfo{
				"pac1": {"1.0", 40},
				"pac2": {"2.0", 50},
				"pac3": {"3.0", 60}},
			map2: map[string]PackageInfo{
				"pac1": {"1.0", 40},
				"pac2": {"2.0", 50},
				"pac3": {"3.0", 60}},
			expected: PackageDiff{
				Packages1: map[string]PackageInfo{},
				Packages2: map[string]PackageInfo{},
				InfoDiff:  []Info{}},
		},
		{
			descrip: "MultiVersion call with identical Packages in different layers",
			map1: map[string]map[string]PackageInfo{
				"pac5": {"globalPath": {"version", 0}},
				"pac3": {"notquite/localPath": {"version", 0}},
				"pac4": {"globalPath": {"version", 0}}},
			map2: map[string]map[string]PackageInfo{
				"pac5": {"globalPath": {"version", 0}},
				"pac3": {"notquite/localPath": {"version", 0}},
				"pac4": {"globalPath": {"version", 0}}},
			expected: MultiVersionPackageDiff{
				Packages1: map[string]map[string]PackageInfo{},
				Packages2: map[string]map[string]PackageInfo{},
				InfoDiff:  []MultiVersionInfo{},
			},
		},
		{
			descrip: "MultiVersion Packages",
			map1: map[string]map[string]PackageInfo{
				"pac5": {"onlyImg1": {"version", 0}},
				"pac4": {"samePlace": {"version", 0}},
				"pac1": {"node_modules/pac1": {"1.0", 40}},
				"pac2": {"usr/local/lib/node_modules/pac2": {"2.0", 50},
					"node_modules/pac2": {"3.0", 50}}},
			map2: map[string]map[string]PackageInfo{
				"pac4": {"samePlace": {"version", 0}},
				"pac1": {"node_modules/pac1": {"2.0", 40}},
				"pac2": {"usr/local/lib/node_modules/pac2": {"4.0", 50}},
				"pac3": {"usr/local/lib/node_modules/pac3": {"5.0", 100}}},
			expected: MultiVersionPackageDiff{
				Packages1: map[string]map[string]PackageInfo{
					"pac5": {"onlyImg1": {"version", 0}},
				},
				Packages2: map[string]map[string]PackageInfo{
					"pac3": {"usr/local/lib/node_modules/pac3": {"5.0", 100}},
				},
				InfoDiff: []MultiVersionInfo{
					{
						Package: "pac1",
						Info1:   []PackageInfo{{"1.0", 40}},
						Info2:   []PackageInfo{{"2.0", 40}},
					},
					{
						Package: "pac2",
						Info1:   []PackageInfo{{"2.0", 50}, {"3.0", 50}},
						Info2:   []PackageInfo{{"4.0", 50}},
					},
				},
			},
		},
	}
	for _, test := range testCases {
		diff := diffMaps(test.map1, test.map2)
		diffVal := reflect.ValueOf(diff)
		testExpVal := reflect.ValueOf(test.expected)
		switch test.expected.(type) {
		case PackageDiff:
			expected := testExpVal.Interface().(PackageDiff)
			actual := diffVal.Interface().(PackageDiff)
			sort.Sort(ByPackage(expected.InfoDiff))
			sort.Sort(ByPackage(actual.InfoDiff))
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected Diff to be: %v but got:%v", expected, actual)
				return
			}
		case MultiVersionPackageDiff:
			expected := testExpVal.Interface().(MultiVersionPackageDiff)
			actual := diffVal.Interface().(MultiVersionPackageDiff)
			sort.Sort(ByMultiPackage(expected.InfoDiff))
			sort.Sort(ByMultiPackage(actual.InfoDiff))
			for _, pack := range expected.InfoDiff {
				sort.Sort(ByPackageInfo(pack.Info1))
				sort.Sort(ByPackageInfo(pack.Info2))
			}
			for _, pack2 := range actual.InfoDiff {
				sort.Sort(ByPackageInfo(pack2.Info1))
				sort.Sort(ByPackageInfo(pack2.Info2))
			}
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected Diff to be: %v but got:%v", expected, actual)
				return
			}
		}
	}
}

func TestCheckPackageMapType(t *testing.T) {
	testCases := []struct {
		descrip       string
		map1          interface{}
		map2          interface{}
		expectedType  reflect.Type
		expectedMulti bool
		err           bool
	}{
		{
			descrip: "Map arguments not maps",
			map1:    "not a map",
			map2:    "not a map either",
			err:     true,
		},
		{
			descrip: "Map arguments not same type",
			map1:    map[string]int{},
			map2:    map[int]string{},
			err:     true,
		},
		{
			descrip:      "Single Version Package Maps",
			map1:         map[string]PackageInfo{},
			map2:         map[string]PackageInfo{},
			expectedType: reflect.TypeOf(map[string]PackageInfo{}),
		},
		{
			descrip:       "MultiVersion Package Maps",
			map1:          map[string]map[string]PackageInfo{},
			map2:          map[string]map[string]PackageInfo{},
			expectedType:  reflect.TypeOf(map[string]map[string]PackageInfo{}),
			expectedMulti: true,
		},
	}
	for _, test := range testCases {
		actualType, actualMulti, err := checkPackageMapType(test.map1, test.map2)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Error("Expected error but got none.")
		}
		if actualType != test.expectedType {
			t.Errorf("Expected type: %s but got: %s", test.expectedType, actualType)
		}
		if actualMulti != test.expectedMulti {
			t.Errorf("Expected multi: %t but got %t", test.expectedMulti, actualMulti)
		}
	}
}
func TestBuildLayerTargets(t *testing.T) {
	testCases := []struct {
		descrip  string
		path     string
		target   string
		expected []string
		err      bool
	}{
		{
			descrip:  "Filter out non directories",
			path:     "testTars/la-croix1-actual",
			target:   "123",
			expected: []string{},
		},
		{
			descrip:  "Error on bad directory path",
			path:     "test_files/notReal",
			target:   "123",
			expected: []string{},
			err:      true,
		},
		{
			descrip:  "Filter out non-directories and get directories",
			path:     "testTars/la-croix3-full",
			target:   "123",
			expected: []string{"testTars/la-croix3-full/nest/123", "testTars/la-croix3-full/nested-dir/123"},
		},
	}
	for _, test := range testCases {
		layers, err := BuildLayerTargets(test.path, test.target)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Errorf("Expected error but got none: %s", err)
		}
		sort.Strings(test.expected)
		sort.Strings(layers)
		if !reflect.DeepEqual(test.expected, layers) {
			t.Errorf("Expected: %s, but got: %s.", test.expected, layers)
		}
	}
}
