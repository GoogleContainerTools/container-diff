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

func TestGetRubyVersion(t *testing.T) {
	testCases := []struct {
		layerPath        string
		expectedVersions []string
		err              bool
	}{
		{
			layerPath:        "testDirs/gemTests/rubyVersionTests/notAFolder",
			expectedVersions: []string{},
			err:              false,
		},
		{
			layerPath:        "testDirs/gemTests/rubyVersionTests/noLibLayer",
			expectedVersions: []string{},
			err:              false,
		},
		{
			layerPath:        "testDirs/gemTests/rubyVersionTests/noRubyLayer",
			expectedVersions: []string{},
			err:              false,
		},
		{
			layerPath:        "testDirs/gemTests/rubyVersionTests/version2.3Layer",
			expectedVersions: []string{"2.3.0"},
			err:              false,
		},
		{
			layerPath:        "testDirs/gemTests/rubyVersionTests/version2.5Layer",
			expectedVersions: []string{"2.5.0"},
			err:              false,
		},
		{
			layerPath:        "testDirs/gemTests/rubyVersionTests/2VersionLayer",
			expectedVersions: []string{"2.5.0", "2.3.0"},
			err:              false,
		},
	}
	for _, test := range testCases {
		version, err := getRubyVersion(test.layerPath)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Error("Expected error but got none.")
		}
		if !reflect.DeepEqual(version, test.expectedVersions) {
			t.Errorf("\nExpected: %s\nGot: %s", test.expectedVersions, version)
		}
	}
}

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
			descrip: "packagesMultiVersion, no GEM_HOME",
			image: pkgutil.Image{
				FSPath: "testDirs/gemTests/packagesMultiVersion",
				Image: &pkgutil.TestImage{
					Config: &v1.ConfigFile{},
				},
			},
			expectedPackages: map[string]map[string]util.PackageInfo{
				"packageone": {
					"/usr/local/lib/ruby/gems/2.3.0/specifications": {Version: "1.0.0", Size: 0},
					"/usr/local/lib/ruby/gems/2.5.0/specifications": {Version: "2.0.0", Size: 0},
					"/usr/lib/ruby/gems/2.3.0/specifications":       {Version: "2.3.0", Size: 0},
				},
				"packagetwo": {"/usr/local/bundle/specifications": {Version: "4.0.1", Size: 0}},
			},
		},
		{
			descrip: "packagesSingleVersion, no GEM_HOME",
			image: pkgutil.Image{
				FSPath: "testDirs/gemTests/packagesSingleVersion",
				Image: &pkgutil.TestImage{
					Config: &v1.ConfigFile{},
				},
			},
			expectedPackages: map[string]map[string]util.PackageInfo{
				"packageone": {"/usr/local/lib/ruby/gems/2.5.0/specifications": {Version: "2.0.0", Size: 0}},
				"packagetwo": {"/usr/local/lib/ruby/gems/2.5.0/specifications": {Version: "4.0.1", Size: 0}},
			},
		},
		{
			descrip: "gemPathTests, GEM_HOME",
			image: pkgutil.Image{
				FSPath: "testDirs/gemTests/gemPathTests",
				Image: &pkgutil.TestImage{
					Config: &v1.ConfigFile{
						Config: v1.Config{
							Env: []string{"GEM_HOME=testDirs/gemTests/gemPathTests/rbenv:/testDirs/gemTests/gemPathTests/usr/local/bundle", "ENVVAR2=something"},
						},
					},
				},
			},
			expectedPackages: map[string]map[string]util.PackageInfo{
				"packageone":   {"/usr/local/lib/ruby/gems/2.5.0/specifications": {Version: "2.0.0", Size: 0}},
				"packagetwo":   {"/usr/local/lib/ruby/gems/2.5.0/specifications": {Version: "4.0.1", Size: 0}},
				"packagethree": {"/rbenv": {Version: "5.0.0", Size: 0}},
				"packagefour":  {"/usr/local/bundle/specifications": {Version: "10.1", Size: 0}},
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
