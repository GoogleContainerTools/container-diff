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

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
	"github.com/google/go-containerregistry/pkg/v1"
)

func TestGetRubyPackages(t *testing.T) {
	testCases := []struct {
		descrip          string
		image            pkgutil.Image
		expectedPackages map[string]map[string]util.PackageInfo
	}{
		{
			descrip: "noPackagesTest",
			image: pkgutil.Image{
				FSPath: "testDirs/gemTests/noPackagesTest",
				Image: &pkgutil.TestImage{
					Config: &v1.ConfigFile{},
				},
			},
			expectedPackages: map[string]map[string]util.PackageInfo{},
		},

		{
			descrip: "packagesMultiVersion",
			image: pkgutil.Image{
				FSPath: "testDirs/gemTests/packagesMultiVersion",
				Image: &pkgutil.TestImage{
					Config: &v1.ConfigFile{},
				},
			},
			expectedPackages: map[string]map[string]util.PackageInfo{
				"packageone": {
					"/usr/local/lib/ruby/gems/2.3.0/specifications/": {Version: "1.0.0", Size: 0},
					"/usr/local/lib/ruby/gems/2.5.0/specifications/": {Version: "2.0.0", Size: 0},
					"/usr/lib/ruby/gems/2.3.0/specifications/":       {Version: "2.3.0", Size: 0},
				},
				"packagetwo": {"/usr/local/bundle/specifications/": {Version: "4.0.1", Size: 0}},
			},
		},
		{
			descrip: "packagesSingleVersion",
			image: pkgutil.Image{
				FSPath: "testDirs/gemTests/packagesSingleVersion",
				Image: &pkgutil.TestImage{
					Config: &v1.ConfigFile{},
				},
			},
			expectedPackages: map[string]map[string]util.PackageInfo{
				"packageone": {"/usr/local/lib/ruby/gems/2.5.0/specifications/": {Version: "2.0.0", Size: 0}},
				"packagetwo": {"/usr/local/lib/ruby/gems/2.5.0/specifications/": {Version: "4.0.1", Size: 0}},
			},
		},
	}
	for _, test := range testCases {
		d := GemAnalyzer{}
		packages, _ := d.getPackages(test.image)
		if !reflect.DeepEqual(packages, test.expectedPackages) {
			t.Errorf("%s\nExpected: %v\nGot: %v", test.descrip, test.expectedPackages, packages)
		}
	}
}
