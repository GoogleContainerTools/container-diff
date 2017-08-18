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
			t.Errorf("\nExpected: %s\nGot: %s", test.expectedVersions, version)
		}
	}
}

func TestGetPythonPackages(t *testing.T) {
	testCases := []struct {
		descrip          string
		image            utils.Image
		expectedPackages map[string]map[string]utils.PackageInfo
	}{
		{
			descrip: "noPackagesTest",
			image: utils.Image{
				FSPath: "testDirs/pipTests/noPackagesTest",
			},
			expectedPackages: map[string]map[string]utils.PackageInfo{},
		},
		{
			descrip: "packagesMultiVersion, no PYTHONPATH",
			image: utils.Image{
				FSPath: "testDirs/pipTests/packagesMultiVersion",
			},
			expectedPackages: map[string]map[string]utils.PackageInfo{
				"packageone": {
					"/usr/local/lib/python3.6/site-packages": {Version: "3.6.9", Size: "0"},
					"/usr/local/lib/python2.7/site-packages": {Version: "0.1.1", Size: "0"},
				},
				"packagetwo": {"/usr/local/lib/python3.6/site-packages": {Version: "4.6.2", Size: "0"}},
				"script1":    {"/usr/local/lib/python3.6/site-packages": {Version: "1.0", Size: "0"}},
				"script2":    {"/usr/local/lib/python3.6/site-packages": {Version: "2.0", Size: "0"}},
				"script3":    {"/usr/local/lib/python2.7/site-packages": {Version: "3.0", Size: "0"}},
			},
		},
		{
			descrip: "packagesSingleVersion, no PYTHONPATH",
			image: utils.Image{
				FSPath: "testDirs/pipTests/packagesSingleVersion",
			},
			expectedPackages: map[string]map[string]utils.PackageInfo{
				"packageone": {"/usr/local/lib/python3.6/site-packages": {Version: "3.6.9", Size: "0"}},
				"packagetwo": {"/usr/local/lib/python3.6/site-packages": {Version: "4.6.2", Size: "0"}},
				"script1":    {"/usr/local/lib/python3.6/site-packages": {Version: "1.0", Size: "0"}},
				"script2":    {"/usr/local/lib/python3.6/site-packages": {Version: "2.0", Size: "0"}},
			},
		},
		{
			descrip: "pythonPathTests, PYTHONPATH",
			image: utils.Image{
				FSPath: "testDirs/pipTests/pythonPathTests",
				Config: utils.ConfigSchema{
					Config: utils.ConfigObject{
						Env: []string{"PYTHONPATH=testDirs/pipTests/pythonPathTests/pythonPath1:testDirs/pipTests/pythonPathTests/pythonPath2/subdir", "ENVVAR2=something"},
					},
				},
			},
			expectedPackages: map[string]map[string]utils.PackageInfo{
				"packageone":   {"/usr/local/lib/python3.6/site-packages": {Version: "3.6.9", Size: "0"}},
				"packagetwo":   {"/usr/local/lib/python3.6/site-packages": {Version: "4.6.2", Size: "0"}},
				"packagefive":  {"/pythonPath2/subdir": {Version: "3.6.9", Size: "0"}},
				"packagesix":   {"/pythonPath1": {Version: "3.6.9", Size: "0"}},
				"packageseven": {"/pythonPath1": {Version: "4.6.2", Size: "0"}},
			},
		},
		{
			descrip: "pythonPathTests, no PYTHONPATH",
			image: utils.Image{
				FSPath: "testDirs/pipTests/pythonPathTests",
				Config: utils.ConfigSchema{
					Config: utils.ConfigObject{
						Env: []string{"ENVVAR=something"},
					},
				},
			},
			expectedPackages: map[string]map[string]utils.PackageInfo{
				"packageone": {"/usr/local/lib/python3.6/site-packages": {Version: "3.6.9", Size: "0"}},
				"packagetwo": {"/usr/local/lib/python3.6/site-packages": {Version: "4.6.2", Size: "0"}},
			},
		},
	}
	for _, test := range testCases {
		d := PipAnalyzer{}
		packages, _ := d.getPackages(test.image)
		if !reflect.DeepEqual(packages, test.expectedPackages) {
			t.Errorf("%s\nExpected: %s\nGot: %s", test.descrip, test.expectedPackages, packages)
		}
	}
}
