package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"code.cloudfoundry.org/bytefmt"
)

// Directory stores a representaiton of a file directory.
type Directory struct {
	Root    string
	Content []string
}

type DirectoryEntry struct {
	Name string
	Size string
}

func GetSize(path string) (int64, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return -1, errors.Errorf("Could not get size for %s: %s", path, err)
	}
	if stat.IsDir() {
		return getDirectorySize(path)
	}
	return stat.Size(), nil
}

func getDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// GetDirectoryContents converts the directory starting at the provided path into a Directory struct.
func GetDirectory(path string, deep bool) (Directory, error) {
	var directory Directory
	directory.Root = path
	var err error
	if deep {
		walkFn := func(currPath string, info os.FileInfo, err error) error {
			newContent := strings.TrimPrefix(currPath, directory.Root)
			if newContent != "" {
				directory.Content = append(directory.Content, newContent)
			}
			return nil
		}

		err = filepath.Walk(path, walkFn)
	} else {
		contents, err := ioutil.ReadDir(path)
		if err != nil {
			return directory, err
		}

		for _, file := range contents {
			fileName := "/" + file.Name()
			directory.Content = append(directory.Content, fileName)
		}
	}
	return directory, err
}

// Checks for content differences between files of the same name from different directories
func GetModifiedEntries(d1, d2 Directory) []string {
	d1files := d1.Content
	d2files := d2.Content

	filematches := GetMatches(d1files, d2files)

	modified := []string{}
	for _, f := range filematches {
		f1path := fmt.Sprintf("%s%s", d1.Root, f)
		f2path := fmt.Sprintf("%s%s", d2.Root, f)

		f1stat, err := os.Stat(f1path)
		if err != nil {
			glog.Errorf("Error checking directory entry %s: %s\n", f, err)
			continue
		}
		f2stat, err := os.Stat(f2path)
		if err != nil {
			glog.Errorf("Error checking directory entry %s: %s\n", f, err)
			continue
		}

		// If the directory entry in question is a tar, verify that the two have the same size
		if isTar(f1path) {
			if f1stat.Size() != f2stat.Size() {
				modified = append(modified, f)
			}
			continue
		}

		// If the directory entry is not a tar and not a directory, then it's a file so make sure the file contents are the same
		// Note: We skip over directory entries because to compare directories, we compare their contents
		if !f1stat.IsDir() {
			same, err := checkSameFile(f1path, f2path)
			if err != nil {
				glog.Errorf("Error diffing contents of %s and %s: %s\n", f1path, f2path, err)
				continue
			}
			if !same {
				modified = append(modified, f)
			}
		}
	}
	return modified
}

func GetAddedEntries(d1, d2 Directory) []string {
	return GetAdditions(d1.Content, d2.Content)
}

func GetDeletedEntries(d1, d2 Directory) []string {
	return GetDeletions(d1.Content, d2.Content)
}

type DirDiff struct {
	Adds   []DirectoryEntry
	Dels   []DirectoryEntry
	Mods   []DirectoryEntry
}

func createDirectoryEntries(entryNames []string) (entries []DirectoryEntry) {
	for _, name := range entryNames {
		strSize := ""
		size, err := GetSize(name)
		if err == nil {
			strSize = bytefmt.ByteSize(size)
		}
		entry := DirectoryEntry{
			Name: name,
			Size: strSize,
		}
		entries = append(entries, entry)
	}
	return entries
}

// DiffDirectory takes the diff of two directories, assuming both are completely unpacked
func DiffDirectory(d1, d2 Directory) (DirDiff, bool) {
	adds := GetAddedEntries(d1, d2)
	dels := GetDeletedEntries(d1, d2)
	mods := GetModifiedEntries(d1, d2)

	var same bool
	if len(adds) == 0 && len(dels) == 0 && len(mods) == 0 {
		same = true
	} else {
		same = false
	}

	return DirDiff{adds, dels, mods}, same
}

func checkSameFile(f1name, f2name string) (bool, error) {
	// Check first if files differ in size and immediately return
	f1stat, err := os.Stat(f1name)
	if err != nil {
		return false, err
	}
	f2stat, err := os.Stat(f2name)
	if err != nil {
		return false, err
	}

	if f1stat.Size() != f2stat.Size() {
		return false, nil
	}

	// Next, check file contents
	f1, err := ioutil.ReadFile(f1name)
	if err != nil {
		return false, err
	}
	f2, err := ioutil.ReadFile(f2name)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(f1, f2) {
		return false, nil
	}
	return true, nil
}
