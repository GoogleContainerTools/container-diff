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
	"reflect"
	"testing"

	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/GoogleCloudPlatform/container-diff/util"
)

func TestGetNodePackages(t *testing.T) {
	testCases := []struct {
		descrip  string
		path     string
		expected map[string]map[string]util.PackageInfo
		err      bool
	}{
		{
			descrip:  "no directory",
			path:     "testDirs/notThere",
			expected: map[string]map[string]util.PackageInfo{},
			err:      true,
		},
		{
			descrip:  "no packages",
			path:     "testDirs/noPackages",
			expected: map[string]map[string]util.PackageInfo{},
		},
		{
			descrip: "all packages in one layer",
			path:    "testDirs/packageOne",
			expected: map[string]map[string]util.PackageInfo{
				"pac1": {"/node_modules/pac1/": {Version: "1.0", Size: 41}},
				"pac2": {"/usr/local/lib/node_modules/pac2/": {Version: "2.0", Size: 41}},
				"pac3": {"/node_modules/pac3/": {Version: "3.0", Size: 41}}},
		},
		{
			descrip: "Multi version packages",
			path:    "testDirs/packageMulti",
			expected: map[string]map[string]util.PackageInfo{
				"pac1": {"/node_modules/pac1/": {Version: "1.0", Size: 41}},
				"pac2": {"/node_modules/pac2/": {Version: "2.0", Size: 41},
					"/usr/local/lib/node_modules/pac2/": {Version: "3.0", Size: 41}}},
		},
	}

	for _, test := range testCases {
		image := pkgutil.Image{FSPath: test.path}
		d := NodeAnalyzer{}
		packages, err := d.getPackages(image)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Errorf("Expected error but got none.")
		}
		if !reflect.DeepEqual(packages, test.expected) {
			t.Errorf("Expected: %v but got: %v", test.expected, packages)
		}
	}
}
func TestReadPackageJSON(t *testing.T) {
	testCases := []struct {
		descrip  string
		path     string
		expected nodePackage
		err      bool
	}{
		{
			descrip: "Error on non-existent file",
			path:    "testDirs/not_real.json",
			err:     true,
		},
		{
			descrip:  "Parse JSON with exact fields",
			path:     "testDirs/exact.json",
			expected: nodePackage{"La-croix", "Lime"},
		},
		{
			descrip:  "Parse JSON with additional fields",
			path:     "testDirs/extra.json",
			expected: nodePackage{"La-croix", "Lime"},
		},
	}
	for _, test := range testCases {
		actual, err := readPackageJSON(test.path)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Error("Expected errorbut got none.")
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected: %s but got: %s", test.expected, actual)
		}
	}
}
