package utils

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
				"pac1": {"1.0", "40"},
				"pac3": {"3.0", "60"}},
			map2: map[string]PackageInfo{
				"pac4": {"4.0", "70"},
				"pac5": {"5.0", "80"}},
			expected: PackageDiff{
				Packages1: map[string]PackageInfo{
					"pac1": {"1.0", "40"},
					"pac3": {"3.0", "60"}},
				Packages2: map[string]PackageInfo{
					"pac4": {"4.0", "70"},
					"pac5": {"5.0", "80"}},
				InfoDiff: []Info{}},
		},
		{
			descrip: "Different Versions and Sizes.",
			map1: map[string]PackageInfo{
				"pac2": {"2.0", "50"},
				"pac3": {"3.0", "60"}},
			map2: map[string]PackageInfo{
				"pac2": {"2.0", "45"},
				"pac3": {"4.0", "60"}},
			expected: PackageDiff{
				Packages1: map[string]PackageInfo{},
				Packages2: map[string]PackageInfo{},
				InfoDiff: []Info{
					{"pac2", PackageInfo{"2.0", "50"}, PackageInfo{"2.0", "45"}},
					{"pac3", PackageInfo{"3.0", "60"}, PackageInfo{"4.0", "60"}}},
			},
		},
		{
			descrip: "Identical packages, versions, and sizes",
			map1: map[string]PackageInfo{
				"pac1": {"1.0", "40"},
				"pac2": {"2.0", "50"},
				"pac3": {"3.0", "60"}},
			map2: map[string]PackageInfo{
				"pac1": {"1.0", "40"},
				"pac2": {"2.0", "50"},
				"pac3": {"3.0", "60"}},
			expected: PackageDiff{
				Packages1: map[string]PackageInfo{},
				Packages2: map[string]PackageInfo{},
				InfoDiff:  []Info{}},
		},
		{
			descrip: "MultiVersion call with identical Packages in different layers",
			map1: map[string]map[string]PackageInfo{
				"pac5": {"img/hash1/globalPath": {"version", "size"}},
				"pac3": {"img/hash1/notquite/localPath": {"version", "size"}},
				"pac4": {"img/samePlace": {"version", "size"}}},
			map2: map[string]map[string]PackageInfo{
				"pac5": {"img/hash2/globalPath": {"version", "size"}},
				"pac3": {"img/hash2/notquite/localPath": {"version", "size"}},
				"pac4": {"img/samePlace": {"version", "size"}}},
			expected: MultiVersionPackageDiff{
				Packages1: map[string]map[string]PackageInfo{},
				Packages2: map[string]map[string]PackageInfo{},
				InfoDiff:  []MultiVersionInfo{},
			},
		},
		{
			descrip: "MultiVersion Packages",
			map1: map[string]map[string]PackageInfo{
				"pac5": {"img/onlyImg1": {"version", "size"}},
				"pac4": {"img/hash1/samePlace": {"version", "size"}},
				"pac1": {"img/layer1/layer/node_modules/pac1": {"1.0", "40"}},
				"pac2": {"img/layer1/layer/usr/local/lib/node_modules/pac2": {"2.0", "50"},
					"img/layer2/layer/usr/local/lib/node_modules/pac2": {"3.0", "50"}}},
			map2: map[string]map[string]PackageInfo{
				"pac4": {"img/hash2/samePlace": {"version", "size"}},
				"pac1": {"img/layer2/layer/node_modules/pac1": {"2.0", "40"}},
				"pac2": {"img/layer3/layer/usr/local/lib/node_modules/pac2": {"4.0", "50"}},
				"pac3": {"img/layer2/layer/usr/local/lib/node_modules/pac2": {"5.0", "100"}}},
			expected: MultiVersionPackageDiff{
				Packages1: map[string]map[string]PackageInfo{
					"pac5": {"img/onlyImg1": {"version", "size"}},
				},
				Packages2: map[string]map[string]PackageInfo{
					"pac3": {"img/layer2/layer/usr/local/lib/node_modules/pac2": {"5.0", "100"}},
				},
				InfoDiff: []MultiVersionInfo{
					{
						Package: "pac1",
						Info1:   []PackageInfo{{"1.0", "40"}},
						Info2:   []PackageInfo{{"2.0", "40"}},
					},
					{
						Package: "pac2",
						Info1:   []PackageInfo{{"2.0", "50"}, {"3.0", "50"}},
						Info2:   []PackageInfo{{"4.0", "50"}},
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
				t.Errorf("expected Diff to be: %s but got:%s", expected, actual)
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
				t.Errorf("expected Diff to be: %s but got:%s", expected, actual)
				return
			}
		}
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		descrip     string
		VersionList []PackageInfo
		Layers      []string
		currLayer   string
		currVersion PackageInfo
		index       int
		ok          bool
	}{
		{
			descrip:     "Does contain",
			VersionList: []PackageInfo{{Version: "2", Size: "b"}, {Version: "1", Size: "a"}},
			Layers:      []string{"img/1/global", "img/2/local"},
			currLayer:   "img/3/local",
			currVersion: PackageInfo{Version: "1", Size: "a"},
			index:       1,
			ok:          true,
		},
		{
			descrip:     "Not contained",
			VersionList: []PackageInfo{{Version: "1", Size: "a"}, {Version: "2", Size: "b"}},
			Layers:      []string{"img/1/global", "img/2/local"},
			currLayer:   "img/3/global",
			currVersion: PackageInfo{Version: "2", Size: "a"},
			index:       0,
			ok:          false,
		},
		{
			descrip:     "Does contain but path doesn't match",
			VersionList: []PackageInfo{{Version: "1", Size: "a"}, {Version: "2", Size: "b"}},
			Layers:      []string{"img/1/local", "img/2/local"},
			currLayer:   "img/3/global",
			currVersion: PackageInfo{Version: "1", Size: "a"},
			index:       0,
			ok:          false,
		},
		{
			descrip:     "Layers and Versions not of same length",
			VersionList: []PackageInfo{{Version: "1", Size: "a"}, {Version: "2", Size: "b"}},
			Layers:      []string{"img/1/local"},
			currLayer:   "img/3/global",
			currVersion: PackageInfo{Version: "1", Size: "a"},
			index:       0,
			ok:          false,
		},
	}
	for _, test := range testCases {
		index, ok := contains(test.VersionList, test.Layers, test.currLayer, test.currVersion)
		if test.ok != ok {
			t.Errorf("Expected status: %t, but got: %t", test.ok, ok)
		}
		if test.index != index {
			t.Errorf("Expected index: %d, but got: %d", test.index, index)
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
