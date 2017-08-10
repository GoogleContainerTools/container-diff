package utils

import (
	"archive/tar"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

// Directory stores a representaiton of a file directory.
type Directory struct {
	Root    string
	Content []string
}

func unpackTar(tr *tar.Reader, path string) error {
	for {
		header, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			glog.Fatalf(err.Error())
			return err
		}

		if strings.Contains(header.Name, ".wh.") {
			rmPath := filepath.Join(path, header.Name)
			newName := strings.Replace(rmPath, ".wh.", "", 1)
			err := os.Remove(rmPath)
			if err != nil {
				glog.Info(err)
			}
			err = os.RemoveAll(newName)
			if err != nil {
				glog.Info(err)
			}
			continue
		}

		target := filepath.Join(path, header.Name)
		mode := header.FileInfo().Mode()
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, mode); err != nil {
					return err
				}
				continue
			}

		// if it's a file create it
		case tar.TypeReg:

			currFile, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
			if err != nil {
				return err
			}
			defer currFile.Close()
			_, err = io.Copy(currFile, tr)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// UnTar takes in a path to a tar file and writes the untarred version to the provided target.
// Only untars one level, does not untar nested tars.
func UnTar(filename string, path string) error {
	if _, ok := os.Stat(path); ok != nil {
		os.MkdirAll(path, 0777)
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	tr := tar.NewReader(file)
	err = unpackTar(tr, path)
	if err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func isTar(path string) bool {
	return filepath.Ext(path) == ".tar"
}

// ExtractTar extracts the tar and any nested tar at the given path.
// After execution the original tar file is removed and the untarred version is in it place.
func ExtractTar(path string) error {
	removeTar := false

	var untarWalkFn func(path string, info os.FileInfo, err error) error

	untarWalkFn = func(path string, info os.FileInfo, err error) error {
		if isTar(path) {
			target := strings.TrimSuffix(path, filepath.Ext(path))
			UnTar(path, target)
			if removeTar {
				os.Remove(path)
			}
			// remove nested tar files that get copied but not the original tar passed
			removeTar = true
			filepath.Walk(target, untarWalkFn)
		}
		return nil
	}

	return filepath.Walk(path, untarWalkFn)
}

func TarToDir(tarPath string, deep bool) (string, string, error) {
	err := ExtractTar(tarPath)
	if err != nil {
		return "", "", err
	}
	path := strings.TrimSuffix(tarPath, filepath.Ext(tarPath))
	jsonPath := path + ".json"
	err = DirToJSON(path, jsonPath, deep)
	if err != nil {
		return "", "", err
	}
	return jsonPath, path, nil
}

// DirToJSON records the directory structure starting at the provided path as in a json file.
func DirToJSON(path string, target string, deep bool) error {
	var directory Directory
	directory.Root = path

	if deep {
		tarJSONWalkFn := func(currPath string, info os.FileInfo, err error) error {
			newContent := strings.TrimPrefix(currPath, directory.Root)
			if newContent != "" {
				directory.Content = append(directory.Content, newContent)
			}
			return nil
		}

		filepath.Walk(path, tarJSONWalkFn)
	} else {
		contents, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}

		for _, file := range contents {
			fileName := "/" + file.Name()
			directory.Content = append(directory.Content, fileName)
		}
	}

	data, err := json.Marshal(directory)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(target, data, 0777)
}

func CheckTar(image string) bool {
	if strings.TrimSuffix(image, ".tar") == image {
		return false
	}
	if _, err := os.Stat(image); err != nil {
		return false
	}
	return true
}
