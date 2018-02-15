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
	"testing"

	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
)

type difftestpair struct {
	input           [2]pkgutil.Directory
	expected_output []string
}

var d1 = pkgutil.Directory{Root: "Dir1", Content: []string{}}
var d2 = pkgutil.Directory{Root: "Dir2", Content: []string{"file1"}}
var d3 = pkgutil.Directory{Root: "Dir2", Content: []string{"file1", "file2"}}

func TestGetAddedEntries(t *testing.T) {
	var additiontests = []difftestpair{
		{[2]pkgutil.Directory{d1, d1}, []string{}},
		{[2]pkgutil.Directory{d2, d1}, []string{}},
		{[2]pkgutil.Directory{d2, d3}, []string{"file2"}},
		{[2]pkgutil.Directory{d1, d3}, []string{"file1", "file2"}},
	}

	for _, test := range additiontests {
		output := GetAddedEntries(test.input[0], test.input[1])
		if !reflect.DeepEqual(output, test.expected_output) {
			t.Errorf("\nExpected: %s\nGot: %s\n", test.expected_output, output)
		}
	}
}

func TestGetDeletedEntries(t *testing.T) {
	var deletiontests = []difftestpair{
		{[2]pkgutil.Directory{d1, d1}, []string{}},
		{[2]pkgutil.Directory{d1, d2}, []string{}},
		{[2]pkgutil.Directory{d3, d2}, []string{"file2"}},
		{[2]pkgutil.Directory{d3, d1}, []string{"file1", "file2"}},
	}

	for _, test := range deletiontests {
		output := GetDeletedEntries(test.input[0], test.input[1])
		if !reflect.DeepEqual(output, test.expected_output) {
			t.Errorf("\nExpected: %s\nGot: %s\n", test.expected_output, output)
		}
	}
}

func TestGetModifiedEntries(t *testing.T) {
	var testdir1 = pkgutil.Directory{Root: "test_files/dir1/", Content: []string{"file1", "file2", "file3"}}
	var testdir2 = pkgutil.Directory{Root: "test_files/dir2/", Content: []string{"file1", "file2", "file4"}}
	var testdir3 = pkgutil.Directory{Root: "test_files/dir1_copy/", Content: []string{"file1", "file2", "file3"}}
	var testdir4 = pkgutil.Directory{Root: "test_files/dir2_modified/", Content: []string{"file1", "file2", "file4"}}

	var modifiedtests = []difftestpair{
		{[2]pkgutil.Directory{d1, d1}, []string{}},
		{[2]pkgutil.Directory{testdir1, testdir3}, []string{}},
		{[2]pkgutil.Directory{testdir1, testdir2}, []string{"file2"}},
		{[2]pkgutil.Directory{testdir2, testdir4}, []string{"file1", "file2", "file4"}},
	}

	for _, test := range modifiedtests {
		output := GetModifiedEntries(test.input[0], test.input[1])
		if !reflect.DeepEqual(output, test.expected_output) {
			t.Errorf("\nExpected: %s\nGot: %s\n", test.expected_output, output)
		}
	}
}

func TestGetDirectory(t *testing.T) {
	tests := []struct {
		descrip  string
		path     string
		expected pkgutil.Directory
		deep     bool
	}{
		{
			descrip: "deep",
			path:    "testTars/la-croix3-full",
			expected: pkgutil.Directory{
				Root:    "testTars/la-croix3-full",
				Content: []string{"/lime.txt", "/nest", "/nest/f1.txt", "/nested-dir", "/nested-dir/f2.txt", "/passionfruit.txt", "/peach-pear.txt"},
			},
			deep: true,
		},
		{
			descrip: "shallow",
			path:    "testTars/la-croix3-full",
			expected: pkgutil.Directory{
				Root:    "testTars/la-croix3-full",
				Content: []string{"/lime.txt", "/nest", "/nested-dir", "/passionfruit.txt", "/peach-pear.txt"},
			},
			deep: false,
		},
	}
	for _, testCase := range tests {
		dir, err := pkgutil.GetDirectory(testCase.path, testCase.deep)
		if err != nil {
			t.Errorf("Error converting directory to Directory struct")
		}

		actualDir := dir
		expectedDir := testCase.expected

		if !reflect.DeepEqual(actualDir, expectedDir) {
			t.Errorf("%s test was incorrect\nExpected: %s\nGot: %s", testCase.descrip, expectedDir, actualDir)
		}
	}
}

func TestCheckSameFile(t *testing.T) {
	type filetestpair struct {
		input            [2]string
		expected_success bool
		expected_output  bool
	}

	var samefiletests = []filetestpair{
		{[2]string{"", ""}, false, false},
		{[2]string{"nonexistent", "file1"}, false, false},
		{[2]string{"test_files/file1", "test_files/file1_copy"}, true, true},
		{[2]string{"test_files/file1", "test_files/file2"}, true, false},
	}

	for _, test := range samefiletests {
		output, err := pkgutil.CheckSameFile(test.input[0], test.input[1])
		if err != nil && test.expected_success {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil {
			if !test.expected_success {
				t.Errorf("Expected error but got none")
			} else {
				if output != test.expected_output {
					t.Errorf("\nExpected: %v\nGot: %v\n", test.expected_output, output)
				}
			}
		}
	}
}

func TestHasFilepathPrefix(t *testing.T) {
	type test struct {
		prefix       string
		path         string
		expectedBool bool
	}

	var tests = []test{
		{"/foo", "/foo/bar", true},
		{"/foo/", "/foo/bar", true},
		{"/foo/bar", "/foo", false},
		{"/foo/bar", "/foo/", false},
		{"/foo", "/foobar", false},
	}

	for _, test := range tests {
		hasPrefix := pkgutil.HasFilepathPrefix(test.path, test.prefix)
		if hasPrefix != test.expectedBool {
			t.Errorf("Got unexpected response: %v for prefix: %s and path: %s", hasPrefix, test.prefix, test.path)
		}
	}
}
