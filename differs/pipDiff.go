package differs

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/container-diff/utils"
	"github.com/golang/glog"
)

type PipDiffer struct {
}

// PipDiff compares pip-installed Python packages between layers of two different images.
func (d PipDiffer) Diff(image1, image2 utils.Image) (utils.DiffResult, error) {
	diff, err := multiVersionDiff(image1, image2, d)
	return diff, err
}

func getPythonVersion(pathToLayer string) ([]string, error) {
	matches := []string{}
	libPath := filepath.Join(pathToLayer, "usr/local/lib")
	libContents, err := ioutil.ReadDir(libPath)
	if err != nil {
		return matches, err
	}

	for _, file := range libContents {
		pattern := regexp.MustCompile("^python[0-9]+\\.[0-9]+$")
		match := pattern.FindString(file.Name())
		if match != "" {
			matches = append(matches, match)
		}
	}
	return matches, nil
}

func getPythonPaths(vars []string) []string {
	paths := []string{}
	for _, envVar := range vars {
		pythonPathPattern := regexp.MustCompile("^PYTHONPATH=(.*)")
		match := pythonPathPattern.FindStringSubmatch(envVar)
		if len(match) != 0 {
			pythonPath := match[1]
			paths = strings.Split(pythonPath, ":")
			break
		}
	}
	return paths
}

func (d PipDiffer) getPackages(image utils.Image) (map[string]map[string]utils.PackageInfo, error) {
	path := image.FSPath
	packages := make(map[string]map[string]utils.PackageInfo)
	pythonPaths := []string{}
	if !reflect.DeepEqual(utils.ConfigSchema{}, image.Config) {
		paths := getPythonPaths(image.Config.Config.Env)
		for _, p := range paths {
			pythonPaths = append(pythonPaths, p)
		}
	}
	pythonVersions, err := getPythonVersion(path)
	if err != nil {
		// Image doesn't have Python installed
		return packages, nil
	}

	for _, pythonVersion := range pythonVersions {
		packagesPath := filepath.Join(path, "usr/local/lib", pythonVersion, "site-packages")
		pythonPaths = append(pythonPaths, packagesPath)
	}

	for _, pythonPath := range pythonPaths {
		contents, err := ioutil.ReadDir(pythonPath)
		if err != nil {
			// python version folder doesn't have a site-packages folder
			continue
		}

		for i := 0; i < len(contents); i++ {
			c := contents[i]
			fileName := c.Name()
			// check if package
			packageDir := regexp.MustCompile("^([a-z|A-Z|0-9|_]+)-(([0-9]+?\\.){2,3})dist-info$")
			packageMatch := packageDir.FindStringSubmatch(fileName)
			if len(packageMatch) != 0 {
				packageName := packageMatch[1]
				version := packageMatch[2][:len(packageMatch[2])-1]

				// Retrieves size for actual package/script corresponding to each dist-info metadata directory
				// by taking the file entry alphabetically before it (for a package) or after it (for a script)
				var size string
				if i-1 >= 0 && contents[i-1].Name() == packageName {
					packagePath := filepath.Join(pythonPath, packageName)
					intSize, err := utils.GetDirectorySize(packagePath)
					if err != nil {
						glog.Errorf("Could not obtain size for package %s", packagePath)
						size = ""
					} else {
						size = strconv.FormatInt(intSize, 10)
					}
				} else if i+1 < len(contents) && contents[i+1].Name() == packageName+".py" {
					size = strconv.FormatInt(contents[i+1].Size(), 10)

				} else {
					glog.Errorf("Could not find Python package %s for corresponding metadata info", packageName)
					continue
				}
				currPackage := utils.PackageInfo{Version: version, Size: size}
				addToMap(packages, packageName, pythonPath, currPackage)
			}
		}
	}

	return packages, nil
}

func addToMap(packages map[string]map[string]utils.PackageInfo, pack string, path string, packInfo utils.PackageInfo) {
	if _, ok := packages[pack]; !ok {
		// package not yet seen
		infoMap := make(map[string]utils.PackageInfo)
		infoMap[path] = packInfo
		packages[pack] = infoMap
		return
	}
	packages[pack][path] = packInfo
}
