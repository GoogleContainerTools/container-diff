/*
Copyright 2020 Google, Inc. All rights reserved.

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

func TestGetEmergePackages(t *testing.T) {
	testCases := []struct {
		descrip  string
		path     string
		expected map[string]util.PackageInfo
		err      bool
	}{
		{
			descrip:  "no directory",
			path:     "testDirs/notThere",
			expected: map[string]util.PackageInfo{},
			err:      true,
		},
		{
			descrip:  "no packages",
			path:     "testDirs/noPackages",
			expected: map[string]util.PackageInfo{},
		},
		{
			descrip: "packages in expected location",
			path:    "testDirs/packageEmerge",
			expected: map[string]util.PackageInfo{
				"dev-python/pkg1": {Version: "0.0.1", Size: 167112},
				"dev-python/pkg2": {Version: "0.0.2", Size: 167112},
				"sys-libs/pkg3":   {Version: "0.0.3", Size: 167112}},
		},
	}
	for _, test := range testCases {
		d := EmergeAnalyzer{}
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
