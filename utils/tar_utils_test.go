package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestUnTar(t *testing.T) {
	testCases := []struct {
		descrip  string
		tarPath  string
		target   string
		expected string
		starter  string
		err      bool
	}{
		{
			descrip:  "Tar with files",
			tarPath:  "testTars/la-croix1.tar",
			target:   "testTars/la-croix1",
			expected: "testTars/la-croix1-actual",
		},
		{
			descrip:  "Tar with folders with files",
			tarPath:  "testTars/la-croix2.tar",
			target:   "testTars/la-croix2",
			expected: "testTars/la-croix2-actual",
		},
		{
			descrip:  "Tar with folders with files and a tar file",
			tarPath:  "testTars/la-croix3.tar",
			target:   "testTars/la-croix3",
			expected: "testTars/la-croix3-actual",
		},
		{
			descrip:  "Tar with .wh.'s",
			tarPath:  "testTars/la-croix-wh.tar",
			target:   "testTars/la-croix-wh",
			expected: "testTars/la-croix-wh-actual",
			starter:  "testTars/la-croix-starter",
		},
		{
			descrip:  "Files updated",
			tarPath:  "testTars/la-croix-update.tar",
			target:   "testTars/la-croix-update",
			expected: "testTars/la-croix-update-actual",
			starter:  "testTars/la-croix-starter",
		},
	}
	for _, test := range testCases {
		if test.starter != "" {
			CopyDir(test.starter, test.target)
		}
		err := UnTar(test.tarPath, test.target)
		if err != nil && !test.err {
			t.Errorf(test.descrip, "Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Errorf(test.descrip, "Expected error but got none: %s", err)
		}
		if !dirEquals(test.expected, test.target) || !dirEquals(test.target, test.expected) {
			t.Errorf(test.descrip, "Directory created not correct structure.")
		}
		os.RemoveAll(test.target)
	}
}

// Copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}

	}

	return nil
}

// Recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.New("Source not a directory")
	}

	// ensure dest dir does not already exist

	_, err = os.Open(dest)
	if !os.IsNotExist(err) {
		return errors.New("Destination already exists")
	}

	// create dest dir

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)

	for _, entry := range entries {

		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			err = CopyDir(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		}

	}
	return nil
}

func TestIsTar(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{input: "/testTar/la-croix1.tar", expected: true},
		{input: "/testTar/la-croix1-actual", expected: false},
	}
	for _, test := range testCases {
		actual := isTar(test.input)
		if test.expected != actual {
			t.Errorf("Expected: %t but got: %t", test.expected, actual)
		}
	}
}

func TestExtractTar(t *testing.T) {
	tarPath := "testTars/la-croix3.tar"
	target := "testTars/la-croix3"
	expected := "testTars/la-croix3-full"
	err := ExtractTar(tarPath)
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
	}
	if !dirEquals(expected, target) || !dirEquals(target, expected) {
		t.Errorf("Directory created not correct structure.")
	}
	os.RemoveAll(target)

}

func dirEquals(actual string, path string) bool {
	files1, _ := ioutil.ReadDir(actual)

	for _, file := range files1 {
		newActualPath := filepath.Join(actual, file.Name())
		newExpectedPath := filepath.Join(path, file.Name())
		fstat, ok := os.Stat(newExpectedPath)
		if ok != nil {
			return false
		}

		if file.IsDir() && !dirEquals(newActualPath, newExpectedPath) {
			return false
		}

		if fstat.Name() != file.Name() {
			return false
		}
		if fstat.Size() != file.Size() {
			return false
		}
		if filepath.Ext(file.Name()) == ".tar" {
			continue
		}

		content1, _ := ioutil.ReadFile(newActualPath)
		content2, _ := ioutil.ReadFile(newExpectedPath)

		if 0 != bytes.Compare(content1, content2) {
			return false
		}
	}
	return true
}

func TestDirToJSON(t *testing.T) {
	tests := []struct {
		path     string
		target   string
		expected string
		deep     bool
	}{
		{"testTars/la-croix3-full", "testTars/la-croix3-full.json", "testTars/la-croix3-actual.json", true},
		{"testTars/la-croix3-full", "testTars/la-croix3-temp.json", "testTars/la-croix3-shallow.json", false},
	}
	for _, testCase := range tests {
		err := DirToJSON(testCase.path, testCase.target, testCase.deep)
		if err != nil {
			t.Errorf("Error converting structure to JSON")
		}

		var actualJSON Directory
		var expectedJSON Directory
		content1, _ := ioutil.ReadFile(testCase.target)
		content2, _ := ioutil.ReadFile(testCase.expected)

		json.Unmarshal(content1, &actualJSON)
		json.Unmarshal(content2, &expectedJSON)

		if !reflect.DeepEqual(actualJSON, expectedJSON) {
			t.Errorf("JSON was incorrect")
		}
		os.Remove(testCase.target)
	}
}
