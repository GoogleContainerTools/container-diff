package cmd

import (
	"bytes"
	"errors"
	goflag "flag"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/GoogleCloudPlatform/container-diff/utils"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var json bool
var eng bool
var save bool
var apt bool
var node bool
var file bool
var history bool
var pip bool

var analyzeFlagMap = map[string]*bool{
	"apt":     &apt,
	"node":    &node,
	"file":    &file,
	"history": &history,
	"pip":     &pip,
}

type validatefxn func(args []string) (bool, error)

var RootCmd = &cobra.Command{
	Use:   "To analyze a single image: [image].  To compare two images: [image1] [image2]",
	Short: "Analyze a single image or compare two images.",
	Long:  `Analyzes a single image or compares two images using the specifed analyzers/differs as indicated via flags (see documentation for available ones).`,
}

func outputResults(resultMap map[string]utils.Result) {
	// Outputs diff/analysis results in alphabetical order by analyzer name
	sortedTypes := []string{}
	for analyzerType := range resultMap {
		sortedTypes = append(sortedTypes, analyzerType)
	}
	sort.Strings(sortedTypes)

	results := make([]interface{}, len(resultMap))
	for i, analyzerType := range sortedTypes {
		result := resultMap[analyzerType]
		if json {
			results[i] = result.OutputStruct()
		} else {
			err := result.OutputText(analyzerType)
			if err != nil {
				glog.Error(err)
			}
		}
	}
	if json {
		err := utils.JSONify(results)
		if err != nil {
			glog.Error(err)
		}
	}
}

func cleanupImage(image utils.Image) {
	if !reflect.DeepEqual(image, (utils.Image{})) {
		glog.Infof("Removing image filesystem directory %s from system", image.FSPath)
		errMsg := remove(image.FSPath, true)
		if errMsg != "" {
			glog.Error(errMsg)
		}
	}
}

func getAllAnalyzers() []string {
	allAnalyzers := []string{}
	for name := range analyzeFlagMap {
		allAnalyzers = append(allAnalyzers, name)
	}
	return allAnalyzers
}

func validateArgs(args []string, validatefxns ...validatefxn) (bool, error) {
	for _, validatefxn := range validatefxns {
		valid, err := validatefxn(args)
		if err != nil {
			return false, err
		} else if !valid {
			return false, nil
		}
	}
	return true, nil
}

func checkImage(arg string) bool {
	if !utils.CheckImageID(arg) && !utils.CheckImageURL(arg) && !utils.CheckTar(arg) {
		return false
	}
	return true
}

func checkArgType(args []string) (bool, error) {
	var buffer bytes.Buffer
	valid := true
	for _, arg := range args {
		if !checkImage(arg) {
			valid = false
			errMessage := fmt.Sprintf("Argument %s is not an image ID, URL, or tar\n", args[0])
			buffer.WriteString(errMessage)
		}
	}
	if !valid {
		return false, errors.New(buffer.String())
	}
	return true, nil
}

func remove(path string, dir bool) string {
	var errStr string
	if path == "" {
		return ""
	}

	var err error
	if dir {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		errStr = "\nUnable to remove " + path
	}
	return errStr
}

func init() {
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func addSharedFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&json, "json", "j", false, "JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).")
	cmd.Flags().BoolVarP(&eng, "eng", "e", false, "By default the docker calls are shelled out locally, set this flag to use the Docker Engine Client (version compatibility required).")
	cmd.Flags().BoolVarP(&pip, "pip", "p", false, "Set this flag to use the pip differ.")
	cmd.Flags().BoolVarP(&node, "node", "n", false, "Set this flag to use the node differ.")
	cmd.Flags().BoolVarP(&apt, "apt", "a", false, "Set this flag to use the apt differ.")
	cmd.Flags().BoolVarP(&file, "file", "f", false, "Set this flag to use the file differ.")
	cmd.Flags().BoolVarP(&history, "history", "d", false, "Set this flag to use the dockerfile history differ.")
	cmd.Flags().BoolVarP(&save, "save", "s", false, "Set this flag to save rather than remove the final image filesystems on exit.")
	cmd.Flags().BoolVarP(&utils.SortSize, "order", "o", false, "Set this flag to sort any file/package results by descending size. Otherwise, they will be sorted by name.")
}
