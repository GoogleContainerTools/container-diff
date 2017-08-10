package utils

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/system"
	"github.com/golang/glog"
)

func GetImageLayers(pathToImage string) []string {
	layers := []string{}
	contents, err := ioutil.ReadDir(pathToImage)
	if err != nil {
		glog.Error(err.Error())
	}

	for _, file := range contents {
		if file.IsDir() {
			layers = append(layers, file.Name())
		}
	}
	return layers
}

func saveImageToTar(image, dest string) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}

	imageTarPath, err := ImageToTar(cli, image, dest)
	if err != nil {
		return "", err
	}
	return imageTarPath, nil
}

// ImageToTar writes an image to a .tar file
func ImageToTar(cli client.APIClient, image, tarName string) (string, error) {
	glog.Info("Saving image")
	imgBytes, err := cli.ImageSave(context.Background(), []string{image})
	if err != nil {
		return "", err
	}
	defer imgBytes.Close()
	newpath := tarName + ".tar"
	return newpath, copyToFile(newpath, imgBytes)
}

func CheckImageID(image string) bool {
	pattern := regexp.MustCompile("[a-z|0-9]{12}")
	if exp := pattern.FindString(image); exp != image {
		return false
	}
	return true
}

func CheckImageURL(image string) bool {
	pattern := regexp.MustCompile("^.+/.+(:.+){0,1}$")
	if exp := pattern.FindString(image); exp != image || CheckTar(image) {
		return false
	}
	return true
}

// copyToFile writes the content of the reader to the specified file
func copyToFile(outfile string, r io.Reader) error {
	// We use sequential file access here to avoid depleting the standby list
	// on Windows. On Linux, this is a call directly to ioutil.TempFile
	tmpFile, err := system.TempFileSequential(filepath.Dir(outfile), ".docker_temp_")
	if err != nil {
		return err
	}

	tmpPath := tmpFile.Name()

	_, err = io.Copy(tmpFile, r)
	tmpFile.Close()

	if err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err = os.Rename(tmpPath, outfile); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}
