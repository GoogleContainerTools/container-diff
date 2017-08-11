package utils

import (
	//"fmt"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
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
		{
			descrip:  "Dir updated",
			tarPath:  "testTars/la-croix-dir-update.tar",
			target:   "testTars/la-croix-dir-update",
			expected: "testTars/la-croix-dir-update-actual",
			starter:  "testTars/la-croix-starter",
		},
	}
	for _, test := range testCases {
		remove := true
		if test.starter != "" {
			CopyDir(test.starter, test.target)
		}
		err := UnTar(test.tarPath, test.target)
		if err != nil && !test.err {
			t.Errorf(test.descrip, "Got unexpected error: %s", err)
			remove = false
		}
		if err == nil && test.err {
			t.Errorf(test.descrip, "Expected error but got none: %s", err)
			remove = false
		}
		if !dirEquals(test.expected, test.target) {
			//d1, _ := GetDirectory(test.expected, true)
			//fmt.Println(d1.Content)
			//d2, _ := GetDirectory(test.target, true)
			//fmt.Println(d2.Content)
			t.Error(test.descrip, ": Directory created not correct structure.")
			remove = false
		}
		if remove {
			os.RemoveAll(test.target)
		}
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

/*func TestExtractTar(t *testing.T) {
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

}*/

func dirEquals(actual string, path string) bool {
	d1, _ := GetDirectory(actual, true)
	d2, _ := GetDirectory(path, true)
	_, same := DiffDirectory(d1, d2)
	return same
}
