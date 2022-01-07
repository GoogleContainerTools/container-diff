/*
Copyright 2018 Google, Inc. All rights reserved.

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

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
)

func TestApkParseLine(t *testing.T) {
	testCases := []struct {
		description string
		line        string
		packages    map[string]util.PackageInfo
		currPackage string
		expPackage  string
		expected    map[string]util.PackageInfo
	}{
		{
			description: "Not applicable line",
			line:        "G:garbage info",
			packages:    map[string]util.PackageInfo{},
			expPackage:  "",
			expected:    map[string]util.PackageInfo{},
		},
		{
			description: "Package line",
			line:        "P:musl",
			currPackage: "foo",
			expPackage:  "musl",
			packages:    map[string]util.PackageInfo{},
			expected:    map[string]util.PackageInfo{},
		},
		{
			description: "Version line",
			line:        "V:1.2.2-r1",
			packages:    map[string]util.PackageInfo{},
			currPackage: "musl",
			expPackage:  "musl",
			expected:    map[string]util.PackageInfo{"musl": {Version: "1.2.2-r1"}},
		},
		{
			description: "Size line",
			line:        "I:622592",
			packages:    map[string]util.PackageInfo{},
			currPackage: "musl",
			expPackage:  "musl",
			expected:    map[string]util.PackageInfo{"musl": {Size: 622592}},
		},
		{
			description: "Pre-existing PackageInfo struct",
			line:        "I:622592",
			packages:    map[string]util.PackageInfo{"musl": {Version: "1.2.2-r1"}},
			currPackage: "musl",
			expPackage:  "musl",
			expected:    map[string]util.PackageInfo{"musl": {Version: "1.2.2-r1", Size: 622592}},
		},
	}

	for _, test := range testCases {
		currPackage := apkParseLine(test.line, test.currPackage, test.packages)
		if currPackage != test.expPackage {
			t.Errorf("Expected current package to be: %s, but got: %s.", test.expPackage, currPackage)
		}
		if !reflect.DeepEqual(test.packages, test.expected) {
			t.Errorf("Expected: %v but got: %v", test.expected, test.packages)
		}
	}
}

func TestGetApkPackages(t *testing.T) {
	testCases := []struct {
		description string
		path        string
		expected    map[string]util.PackageInfo
		err         bool
	}{
		{
			description: "no directory",
			path:        "testDirs/notThere",
			expected:    map[string]util.PackageInfo{},
			err:         true,
		},
		{
			description: "no packages",
			path:        "testDirs/noPackages",
			expected:    map[string]util.PackageInfo{},
		},
		{
			description: "packages in expected location",
			path:        "testDirs/packageOne",
			expected: map[string]util.PackageInfo{
				"musl":    {Version: "1.2.2-r1", Size: 622592},
				"busybox": {Version: "1.32.1-r7", Size: 946176}},
		},
	}
	for _, test := range testCases {
		d := ApkAnalyzer{}
		image := pkgutil.Image{FSPath: test.path}
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
