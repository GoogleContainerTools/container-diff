package utils

import (
	"reflect"
	"testing"
)

type difftestpair struct {
	input           [2]Directory
	expected_output []string
}

var d1 = Directory{"Dir1", []string{}}
var d2 = Directory{"Dir2", []string{"file1"}}
var d3 = Directory{"Dir2", []string{"file1", "file2"}}

func TestGetAddedEntries(t *testing.T) {
	var additiontests = []difftestpair{
		{[2]Directory{d1, d1}, []string{}},
		{[2]Directory{d2, d1}, []string{}},
		{[2]Directory{d2, d3}, []string{"file2"}},
		{[2]Directory{d1, d3}, []string{"file1", "file2"}},
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
		{[2]Directory{d1, d1}, []string{}},
		{[2]Directory{d1, d2}, []string{}},
		{[2]Directory{d3, d2}, []string{"file2"}},
		{[2]Directory{d3, d1}, []string{"file1", "file2"}},
	}

	for _, test := range deletiontests {
		output := GetDeletedEntries(test.input[0], test.input[1])
		if !reflect.DeepEqual(output, test.expected_output) {
			t.Errorf("\nExpected: %s\nGot: %s\n", test.expected_output, output)
		}
	}
}

func TestGetModifiedEntries(t *testing.T) {
	var testdir1 = Directory{"test_files/dir1/", []string{"file1", "file2", "file3"}}
	var testdir2 = Directory{"test_files/dir2/", []string{"file1", "file2", "file4"}}
	var testdir3 = Directory{"test_files/dir1_copy/", []string{"file1", "file2", "file3"}}
	var testdir4 = Directory{"test_files/dir2_modified/", []string{"file1", "file2", "file4"}}

	var modifiedtests = []difftestpair{
		{[2]Directory{d1, d1}, []string{}},
		{[2]Directory{testdir1, testdir3}, []string{}},
		{[2]Directory{testdir1, testdir2}, []string{"file2"}},
		{[2]Directory{testdir2, testdir4}, []string{"file1", "file2", "file4"}},
	}

	for _, test := range modifiedtests {
		output := GetModifiedEntries(test.input[0], test.input[1])
		if !reflect.DeepEqual(output, test.expected_output) {
			t.Errorf("\nExpected: %s\nGot: %s\n", test.expected_output, output)
		}
	}
}

func TestGetDirectory(t *testing.T) {
	type dirtestpair struct {
		input            string
		expected_success bool
	}

	var dirtests = []dirtestpair{
		{"test_files/dir.json", true},
		{"test_files/dir_bad.json", false},
		{"nonexistentpath", false},
		{"", false},
	}

	for _, test := range dirtests {
		_, err := GetDirectory(test.input)
		if err != nil && test.expected_success {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && !test.expected_success {
			t.Errorf("Expected error but got none")
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
		output, err := checkSameFile(test.input[0], test.input[1])
		if err != nil && test.expected_success {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil {
			if !test.expected_success {
				t.Errorf("Expected error but got none")
			} else {
				if output != test.expected_output {
					t.Errorf("\nExpected: %s\nGot: %s\n", test.expected_output, output)
				}
			}
		}
	}
}
