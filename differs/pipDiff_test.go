package differs

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/container-diff/utils"
)

func TestGetPythonVersion(t *testing.T) {
	testCases := []struct {
		layerPath        string
		expectedVersions []string
		err              bool
	}{
		{
			layerPath:        "testDirs/pipTests/pythonVersionTests/notAFolder",
			expectedVersions: []string{},
			err:              true,
		},
		{
			layerPath:        "testDirs/pipTests/pythonVersionTests/noLibLayer",
			expectedVersions: []string{},
			err:              true,
		},
		{
			layerPath:        "testDirs/pipTests/pythonVersionTests/noPythonLayer",
			expectedVersions: []string{},
			err:              false,
		},
		{
			layerPath:        "testDirs/pipTests/pythonVersionTests/version2.7Layer",
			expectedVersions: []string{"python2.7"},
			err:              false,
		},
		{
			layerPath:        "testDirs/pipTests/pythonVersionTests/version3.6Layer",
			expectedVersions: []string{"python3.6"},
			err:              false,
		},
		{
			layerPath:        "testDirs/pipTests/pythonVersionTests/2VersionLayer",
			expectedVersions: []string{"python2.7", "python3.6"},
			err:              false,
		},
	}
	for _, test := range testCases {
		version, err := getPythonVersion(test.layerPath)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Error("Expected error but got none.")
		}
		if !reflect.DeepEqual(version, test.expectedVersions) {
			t.Errorf("Expected: %s.  Got: %s", test.expectedVersions, version)
		}
	}
}

func TestGetPythonPackages(t *testing.T) {
	testCases := []struct {
		path             string
		expectedPackages map[string]map[string]utils.PackageInfo
	}{
		{
			path:             "testDirs/pipTests/noPackagesTest",
			expectedPackages: map[string]map[string]utils.PackageInfo{},
		},
		{
			path: "testDirs/pipTests/packagesOneLayer",
			expectedPackages: map[string]map[string]utils.PackageInfo{
				"packageone": {"python3.6": {Version: "3.6.9", Size: "0"}},
				"packagetwo": {"python3.6": {Version: "4.6.2", Size: "0"}},
				"script1":    {"python3.6": {Version: "1.0", Size: "0"}},
				"script2":    {"python3.6": {Version: "2.0", Size: "0"}},
			},
		},
		{
			path: "testDirs/pipTests/packagesMultiVersion",
			expectedPackages: map[string]map[string]utils.PackageInfo{
				"packageone": {"python3.6": {Version: "3.6.9", Size: "0"},
					"python2.7": {Version: "0.1.1", Size: "0"}},
				"packagetwo": {"python3.6": {Version: "4.6.2", Size: "0"}},
				"script1":    {"python3.6": {Version: "1.0", Size: "0"}},
				"script2":    {"python3.6": {Version: "2.0", Size: "0"}},
				"script3":    {"python2.7": {Version: "3.0", Size: "0"}},
			},
		},
	}
	for _, test := range testCases {
		d := PipDiffer{}
		packages, _ := d.getPackages(test.path)
		if !reflect.DeepEqual(packages, test.expectedPackages) {
			t.Errorf("Expected: %s but got: %s", test.expectedPackages, packages)
		}
	}
}
