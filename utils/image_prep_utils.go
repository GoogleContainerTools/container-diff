package utils

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/containers/image/docker"
	"github.com/docker/docker/api/types/container"
	"github.com/golang/glog"
)

var sourceToPrepMap = map[string]Prepper{
	"ID":  IDPrepper{},
	"URL": CloudPrepper{},
	"tar": TarPrepper{},
}

var sourceCheckMap = map[string]func(string) bool{
	"ID":  CheckImageID,
	"URL": CheckImageURL,
	"tar": CheckTar,
}

type Image struct {
	Source string
	FSPath string
	Config ConfigSchema
}

type ImageHistoryItem struct {
	CreatedBy string `json:"created_by"`
}

type ConfigObject struct {
	Env []string `json:"Env"`
}

type ConfigSchema struct {
	Config  ConfigObject       `json:"config"`
	History []ImageHistoryItem `json:"history"`
}

type ImagePrepper struct {
	Source string
}

type Prepper interface {
	getFileSystem() (string, error)
	getConfig() (ConfigSchema, error)
}

func (p ImagePrepper) GetImage() (Image, error) {
	glog.Infof("Starting prep for image %s", p.Source)
	img := p.Source

	var prepper Prepper
	for source, check := range sourceCheckMap {
		if check(img) {
			typePrepper := reflect.TypeOf(sourceToPrepMap[source])
			prepper = reflect.New(typePrepper).Interface().(Prepper)
			reflect.ValueOf(prepper).Elem().Field(0).Set(reflect.ValueOf(p))
			break
		}
	}
	if prepper == nil {
		return Image{}, errors.New("Could not retrieve image from source")
	}

	imgPath, err := prepper.getFileSystem()
	if err != nil {
		return Image{}, err
	}

	config, err := prepper.getConfig()
	if err != nil {
		glog.Error("Error retrieving History: ", err)
	}

	glog.Infof("Finished prepping image %s", p.Source)
	return Image{
		Source: img,
		FSPath: imgPath,
		Config: config,
	}, nil
}

func getImageFromTar(tarPath string) (string, error) {
	glog.Info("Extracting image tar to obtain image file system")
	path := strings.TrimSuffix(tarPath, filepath.Ext(tarPath))
	err := UnTar(tarPath, path)
	return path, err
}

// CloudPrepper prepares images sourced from a Cloud registry
type CloudPrepper struct {
	ImagePrepper
}

func (p CloudPrepper) getFileSystem() (string, error) {
	// The regexp when passed a string creates a list of the form
	// [repourl/image:tag, image:tag, tag] (the tag may or may not be present)
	URLPattern := regexp.MustCompile("^.+/(.+(:.+){0,1})$")
	URLMatch := URLPattern.FindStringSubmatch(p.Source)
	// Removing the ":" so that the image path name can be <image><tag>
	path := strings.Replace(URLMatch[1], ":", "", -1)
	ref, err := docker.ParseReference("//" + p.Source)
	if err != nil {
		panic(err)
	}

	img, err := ref.NewImage(nil)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	defer img.Close()

	imgSrc, err := ref.NewImageSource(nil, nil)
	if err != nil {
		glog.Error(err)
		return "", err
	}

	if _, ok := os.Stat(path); ok != nil {
		os.MkdirAll(path, 0777)
	}

	for _, b := range img.LayerInfos() {
		bi, _, err := imgSrc.GetBlob(b)
		if err != nil {
			glog.Errorf("Diff may be inaccurate, failed to pull image layer with error: %s", err)
		}
		gzf, err := gzip.NewReader(bi)
		if err != nil {
			glog.Errorf("Diff may be inaccurate, failed to read layers with error: %s", err)
		}
		tr := tar.NewReader(gzf)
		err = unpackTar(tr, path)
		if err != nil {
			glog.Errorf("Diff may be inaccurate, failed to untar layer with error: %s", err)
		}
	}
	return path, nil
}

func (p CloudPrepper) getConfig() (ConfigSchema, error) {
	ref, err := docker.ParseReference("//" + p.Source)
	if err != nil {
		return ConfigSchema{}, err
	}

	img, err := ref.NewImage(nil)
	if err != nil {
		glog.Errorf("Error referencing image %s from registry: %s", p.Source, err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}
	defer img.Close()

	configBlob, err := img.ConfigBlob()
	if err != nil {
		glog.Errorf("Error obtaining config blob for image %s from registry: %s", p.Source, err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}

	var config ConfigSchema
	err = json.Unmarshal(configBlob, &config)
	if err != nil {
		glog.Errorf("Error with config file struct for image %s: %s", p.Source, err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}
	return config, nil
}

type IDPrepper struct {
	ImagePrepper
}

func (p IDPrepper) getFileSystem() (string, error) {
	// check client compatibility with Docker API
	valid, err := ValidDockerVersion()
	if err != nil {
		return "", err
	}
	var tarPath string
	if !valid {
		glog.Info("Docker version incompatible with api, shelling out to local Docker client.")
		tarPath, err = imageToTarCmd(p.Source, p.Source)
	} else {
		tarPath, err = saveImageToTar(p.Source, p.Source)
	}
	if err != nil {
		return "", err
	}

	defer os.Remove(tarPath)
	glog.Info("Extracting image tar to obtain image file system")
	path := strings.TrimSuffix(tarPath, filepath.Ext(tarPath))
	err = ExtractTar(tarPath)
	return path, err
}

func (p IDPrepper) getConfig() (ConfigSchema, error) {
	// check client compatibility with Docker API
	valid, err := ValidDockerVersion()
	if err != nil {
		return ConfigSchema{}, err
	}
	var containerConfig container.Config
	if !valid {
		glog.Info("Docker version incompatible with api, shelling out to local Docker client.")
		containerConfig, err = getImageConfigCmd(p.Source)
	} else {
		containerConfig, err = getImageConfig(p.Source)
	}
	if err != nil {
		return ConfigSchema{}, err
	}

	config := ConfigObject{
		Env: containerConfig.Env,
	}
	history := p.getHistory()
	return ConfigSchema{
		Config:  config,
		History: history,
	}, nil
}

func (p IDPrepper) getHistory() []ImageHistoryItem {
	history, err := getImageHistory(p.Source)
	if err != nil {
		glog.Error("Could not obtain image history for %s: %s", p.Source, err)
	}
	historyItems := []ImageHistoryItem{}
	for _, item := range history {
		historyItems = append(historyItems, ImageHistoryItem{CreatedBy: item.CreatedBy})
	}
	return historyItems
}

type TarPrepper struct {
	ImagePrepper
}

func (p TarPrepper) getFileSystem() (string, error) {
	return getImageFromTar(p.Source)
}

func (p TarPrepper) getConfig() (ConfigSchema, error) {
	tmpDir := strings.TrimSuffix(p.Source, filepath.Ext(p.Source))
	defer os.Remove(tmpDir)
	err := UnTar(p.Source, tmpDir)
	if err != nil {
		return ConfigSchema{}, err
	}
	contents, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		glog.Errorf("Could not read image tar contents: %s", err)
		return ConfigSchema{}, errors.New("Could not obtain image config")
	}

	var config ConfigSchema
	configList := []string{}
	for _, item := range contents {
		if filepath.Ext(item.Name()) == ".json" && item.Name() != "manifest.json" {
			if len(configList) != 0 {
				// Another <image>.json file has already been processed and the image config obtained is uncertain.
				glog.Error("Multiple possible config sources detected for image at " + p.Source + ". Multiple images likely contained in tar. Choosing first one, but diff results may not be completely accurate.")
				break
			}
			fileName := filepath.Join(tmpDir, item.Name())
			file, err := ioutil.ReadFile(fileName)
			if err != nil {
				glog.Errorf("Could not read config file %s: %s", fileName, err)
				return ConfigSchema{}, errors.New("Could not obtain image config")
			}
			err = json.Unmarshal(file, &config)
			if err != nil {
				glog.Errorf("Could not marshal config file %s: %s", fileName, err)
				return ConfigSchema{}, errors.New("Could not obtain image config")
			}
			configList = append(configList, fileName)
		}
	}
	if reflect.DeepEqual(ConfigSchema{}, config) {
		glog.Warningf("No image config found in tar source %s. Pip differ may be incomplete due to missing PYTHONPATH information.")
		return config, nil
	}
	return config, nil
}
