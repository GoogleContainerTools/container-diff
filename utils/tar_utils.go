package utils

import (
	"archive/tar"
	"io"
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
func UnTar(filename string, target string) error {
	if _, ok := os.Stat(target); ok != nil {
		os.MkdirAll(target, 0777)
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	tr := tar.NewReader(file)
	err = unpackTar(tr, target)
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

func CheckTar(image string) bool {
	if strings.TrimSuffix(image, ".tar") == image {
		return false
	}
	if _, err := os.Stat(image); err != nil {
		glog.Errorf("%s does not exist", image)
		return false
	}
	return true
}
