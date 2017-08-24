package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types/container"
	img "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/golang/glog"
)

var eng bool

func SetDockerEngine(useDocker bool) {
	eng = useDocker
}

// ValidDockerVersion determines if there is a Docker client of the necessary version locally installed.
func ValidDockerVersion() (bool, error) {
	_, err := client.NewEnvClient()
	if err != nil {
		return false, fmt.Errorf("Docker client error: %s", err)
	}
	if eng {
		return true, nil
	}
	return false, nil
}

type HistDiff struct {
	Adds []string
	Dels []string
}

// getImageHistory shells out the docker history command and returns a list of history response items.
// The history response items contain only the Created By information for each event.
func getImageHistoryCmd(image string) ([]img.HistoryResponseItem, error) {
	imageID := image
	var err error
	var history []img.HistoryResponseItem
	histArgs := []string{"history", "--no-trunc", imageID}
	dockerHistCmd := exec.Command("docker", histArgs...)
	var response bytes.Buffer
	dockerHistCmd.Stdout = &response
	if err := dockerHistCmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				glog.Error("Docker History Command Exit Status: ", status.ExitStatus())
			}
		} else {
			return history, err
		}
	}
	history, err = processHistOutput(response)
	if err != nil {
		return history, err
	}
	return history, nil

}

func processHistOutput(response bytes.Buffer) ([]img.HistoryResponseItem, error) {
	respReader := bytes.NewReader(response.Bytes())
	reader := bufio.NewReader(respReader)
	var history []img.HistoryResponseItem
	var CreatedByIndex int
	var SizeIndex int
	for {
		var event img.HistoryResponseItem
		text, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return history, err
		}

		line := string(text)
		if CreatedByIndex == 0 {
			CreatedByIndex = strings.Index(line, "CREATED BY")
			SizeIndex = strings.Index(line, "SIZE")
			continue
		}
		event.CreatedBy = line[CreatedByIndex:SizeIndex]
		history = append(history, event)
	}
	return history, nil
}

func imageToTar(image, dest string) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}

	imageTarPath, err := saveImageToTar(cli, image, dest)
	if err != nil {
		return "", err
	}
	return imageTarPath, nil
}

// ImageToTar writes an image to a .tar file
func saveImageToTar(cli client.APIClient, image, tarName string) (string, error) {
	glog.Info("Saving image")
	imgBytes, err := cli.ImageSave(context.Background(), []string{image})
	if err != nil {
		return "", err
	}
	defer imgBytes.Close()
	newpath := tarName + ".tar"
	return newpath, copyToFile(newpath, imgBytes)
}

func imageToTarCmd(imageID, imageName string) (string, error) {
	glog.Info("Saving image")
	cmdArgs := []string{"save", imageID}
	dockerSaveCmd := exec.Command("docker", cmdArgs...)
	var out bytes.Buffer
	dockerSaveCmd.Stdout = &out
	if err := dockerSaveCmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				glog.Error("Docker Save Command Exit Status: ", status.ExitStatus())
			}
		} else {
			return "", err
		}
	}
	imageTarPath := imageName + ".tar"
	reader := bytes.NewReader(out.Bytes())
	err := copyToFile(imageTarPath, reader)
	if err != nil {
		return "", err
	}
	return imageTarPath, nil
}

func getImageHistory(image string) ([]img.HistoryResponseItem, error) {
	validDocker, err := ValidDockerVersion()
	if err != nil {
		return []img.HistoryResponseItem{}, err
	}
	var history []img.HistoryResponseItem
	if validDocker {
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			return []img.HistoryResponseItem{}, err
		}
		history, err = cli.ImageHistory(ctx, image)
		if err != nil {
			return []img.HistoryResponseItem{}, err
		}
	} else {
		glog.Info("Docker version incompatible with api, shelling out to local Docker client.")
		history, err = getImageHistoryCmd(image)
		if err != nil {
			return []img.HistoryResponseItem{}, err
		}
	}
	return history, nil
}

func getImageConfigCmd(image string) (container.Config, error) {
	var err error
	var config container.Config
	configArgs := []string{"inspect", "--format='{{json .Config}}'", image}
	dockerInspectCmd := exec.Command("docker", configArgs...)
	var response bytes.Buffer
	dockerInspectCmd.Stdout = &response
	if err := dockerInspectCmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				glog.Error("Docker Inspect Command Exit Status: ", status.ExitStatus())
			}
		} else {
			return config, err
		}
	}
	err = json.Unmarshal(response.Bytes(), &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func getImageConfig(image string) (container.Config, error) {
	validDocker, err := ValidDockerVersion()
	if err != nil {
		return container.Config{}, err
	}
	var config container.Config
	if validDocker {
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			return container.Config{}, err
		}
		inspect, _, err := cli.ImageInspectWithRaw(ctx, image)
		if err != nil {
			return container.Config{}, err
		}
		config = *inspect.Config
	} else {
		glog.Info("Docker version incompatible with api, shelling out to local Docker client.")
		config, err = getImageConfigCmd(image)
		if err != nil {
			return container.Config{}, err
		}
	}
	return config, nil
}

func getLayersFromManifest(manifestPath string) ([]string, error) {
	type Manifest struct {
		Layers []string
	}

	manifestJSON, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		errMsg := fmt.Sprintf("Could not open manifest to get layer order: %s", err)
		return []string{}, errors.New(errMsg)
	}

	var imageManifest []Manifest
	err = json.Unmarshal(manifestJSON, &imageManifest)
	if err != nil {
		errMsg := fmt.Sprintf("Could not unmarshal manifest to get layer order: %s", err)
		return []string{}, errors.New(errMsg)
	}
	return imageManifest[0].Layers, nil
}

func unpackDockerSave(tarPath string, target string) error {
	if _, ok := os.Stat(target); ok != nil {
		os.MkdirAll(target, 0777)
	}

	tempLayerDir := target + "-temp"
	err := UnTar(tarPath, tempLayerDir)
	if err != nil {
		errMsg := fmt.Sprintf("Could not unpack saved Docker image %s: %s", tarPath, err)
		return errors.New(errMsg)
	}

	manifest := filepath.Join(tempLayerDir, "manifest.json")
	layers, err := getLayersFromManifest(manifest)
	if err != nil {
		return err
	}

	for _, layer := range layers {
		layerTar := filepath.Join(tempLayerDir, layer)
		if _, err := os.Stat(layerTar); err != nil {
			glog.Infof("Did not unpack layer %s because no layer.tar found", layer)
			continue
		}
		err = UnTar(layerTar, target)
		if err != nil {
			glog.Errorf("Could not unpack layer %s: %s", layer, err)
		}
	}
	err = os.RemoveAll(tempLayerDir)
	if err != nil {
		glog.Errorf("Error deleting temp image layer filesystem: %s", err)
	}
	return nil
}
