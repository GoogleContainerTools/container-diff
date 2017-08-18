package differs

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/container-diff/utils"
)

func TestParseLine(t *testing.T) {
	testCases := []struct {
		descrip     string
		line        string
		packages    map[string]utils.PackageInfo
		currPackage string
		expPackage  string
		expected    map[string]utils.PackageInfo
	}{
		{
			descrip:    "Not applicable line",
			line:       "Garbage: garbage info",
			packages:   map[string]utils.PackageInfo{},
			expPackage: "",
			expected:   map[string]utils.PackageInfo{},
		},
		{
			descrip:     "Package line",
			line:        "Package: La-Croix",
			currPackage: "Tea",
			expPackage:  "La-Croix",
			packages:    map[string]utils.PackageInfo{},
			expected:    map[string]utils.PackageInfo{},
		},
		{
			descrip:     "Version line",
			line:        "Version: Lime",
			packages:    map[string]utils.PackageInfo{},
			currPackage: "La-Croix",
			expPackage:  "La-Croix",
			expected:    map[string]utils.PackageInfo{"La-Croix": {Version: "Lime"}},
		},
		{
			descrip:     "Version line with deb release info",
			line:        "Version: Lime+extra_lime",
			packages:    map[string]utils.PackageInfo{},
			currPackage: "La-Croix",
			expPackage:  "La-Croix",
			expected:    map[string]utils.PackageInfo{"La-Croix": {Version: "Lime extra_lime"}},
		},
		{
			descrip:     "Size line",
			line:        "Installed-Size: 12floz",
			packages:    map[string]utils.PackageInfo{},
			currPackage: "La-Croix",
			expPackage:  "La-Croix",
			expected:    map[string]utils.PackageInfo{"La-Croix": {Size: "12floz"}},
		},
		{
			descrip:     "Pre-existing PackageInfo struct",
			line:        "Installed-Size: 12floz",
			packages:    map[string]utils.PackageInfo{"La-Croix": {Version: "Lime"}},
			currPackage: "La-Croix",
			expPackage:  "La-Croix",
			expected:    map[string]utils.PackageInfo{"La-Croix": {Version: "Lime", Size: "12floz"}},
		},
	}

	for _, test := range testCases {
		currPackage := parseLine(test.line, test.currPackage, test.packages)
		if currPackage != test.expPackage {
			t.Errorf("Expected current package to be: %s, but got: %s.", test.expPackage, currPackage)
		}
		if !reflect.DeepEqual(test.packages, test.expected) {
			t.Errorf("Expected: %s but got: %s", test.expected, test.packages)
		}
	}
}

func TestGetAptPackages(t *testing.T) {
	testCases := []struct {
		descrip  string
		path     string
		expected map[string]utils.PackageInfo
		err      bool
	}{
		{
			descrip:  "no directory",
			path:     "testDirs/notThere",
			expected: map[string]utils.PackageInfo{},
			err:      true,
		},
		{
			descrip:  "no packages",
			path:     "testDirs/noPackages",
			expected: map[string]utils.PackageInfo{},
		},
		{
			descrip: "packages in expected location",
			path:    "testDirs/packageOne",
			expected: map[string]utils.PackageInfo{
				"pac1": {Version: "1.0"},
				"pac2": {Version: "2.0"},
				"pac3": {Version: "3.0"}},
		},
	}
	for _, test := range testCases {
		d := AptAnalyzer{}
		image := utils.Image{FSPath: test.path}
		packages, err := d.getPackages(image)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Errorf("Expected error but got none.")
		}
		if !reflect.DeepEqual(packages, test.expected) {
			t.Errorf("Expected: %s but got: %s", test.expected, packages)
		}
	}
}
