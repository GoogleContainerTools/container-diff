package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types"
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
func getImagePullResponse(image string, response []Event) (string, error) {
	var imageDigest string
	for _, event := range response {
		if event.Error != "" {
			err := fmt.Errorf("Error pulling image %s: %s", image, event.Error)
			return "", err
		}
		digestPattern := regexp.MustCompile("^Digest: (sha256:[a-z|0-9]{64})$")
		digestMatch := digestPattern.FindStringSubmatch(event.Status)
		if len(digestMatch) != 0 {
			imageDigest = digestMatch[1]
			return imageDigest, nil
		}
	}
	err := fmt.Errorf("Could not pull image %s", image)
	return "", err
}

func processImagePullEvents(image string, events []Event) (string, string, error) {
	imageDigest, err := getImagePullResponse(image, events)
	if err != nil {
		return "", "", err
	}

	URLPattern := regexp.MustCompile("^.+/(.+(:.+){0,1})$")
	URLMatch := URLPattern.FindStringSubmatch(image)
	imageName := strings.Replace(URLMatch[1], ":", "", -1)
	imageURL := strings.TrimSuffix(image, URLMatch[2])
	imageID := imageURL + "@" + imageDigest

	return imageID, imageName, nil
}

type Event struct {
	Status         string `json:"status"`
	Error          string `json:"error"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func pullImageFromRepo(image string) (string, string, error) {
	glog.Info("Pulling image")
	cli, err := client.NewEnvClient()
	response, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return "", "", err
	}
	defer response.Close()

	d := json.NewDecoder(response)

	var events []Event
	for {
		var event Event
		if err := d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			return "", "", err
		}
		events = append(events, event)
	}
	return processImagePullEvents(image, events)
}

type HistDiff struct {
	Image1 string
	Image2 string
	Adds   []string
	Dels   []string
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

func processPullCmdOutput(image string, response bytes.Buffer) (string, string, error) {
	respReader := bytes.NewReader(response.Bytes())
	reader := bufio.NewReader(respReader)

	var events []Event
	for {
		var event Event
		text, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", "", err
		}
		event.Status = string(text)
		events = append(events, event)
	}
	return processImagePullEvents(image, events)
}

func pullImageCmd(image string) (string, string, error) {
	glog.Info("Pulling image")
	pullArgs := []string{"pull", image}
	dockerPullCmd := exec.Command("docker", pullArgs...)
	var response bytes.Buffer
	dockerPullCmd.Stdout = &response
	if err := dockerPullCmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				glog.Error("Docker Pull Command Exit Status: ", status.ExitStatus())
			}
		} else {
			return "", "", err
		}
	}
	return processPullCmdOutput(image, response)
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
