package cmd

import (
	"bytes"
	"errors"
	goflag "flag"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/GoogleCloudPlatform/container-diff/differs"
	"github.com/GoogleCloudPlatform/container-diff/utils"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var json bool
var eng bool

var apt bool
var node bool
var file bool
var history bool
var pip bool

var save bool

var diffFlagMap = map[string]*bool{
	"apt":     &apt,
	"node":    &node,
	"file":    &file,
	"history": &history,
	"pip":     &pip,
}

var RootCmd = &cobra.Command{
	Use:   "[image1] [image2]",
	Short: "Compare two images.",
	Long:  `Compares two images using the specifed differs as indicated via flags (see documentation for available differs).`,
	Run: func(cmd *cobra.Command, args []string) {
		if validArgs, err := validateArgs(args); !validArgs {
			glog.Error(err.Error())
			os.Exit(1)
		}

		utils.SetDockerEngine(eng)

		img1Arg := args[0]
		img2Arg := args[1]
		diffArgs := []string{}
		allDiffers := getAllDiffers()
		for _, name := range allDiffers {
			if *diffFlagMap[name] == true {
				diffArgs = append(diffArgs, name)
			}
		}
		// If no differs are specified, perform all diffs as the default
		if len(diffArgs) == 0 {
			diffArgs = allDiffers
		}

		var wg sync.WaitGroup
		wg.Add(2)

		glog.Infof("Starting diff on images %s and %s, using differs: %s", img1Arg, img2Arg, diffArgs)

		var image1, image2 utils.Image
		var err error
		go func() {
			defer wg.Done()
			image1, err = utils.ImagePrepper{img1Arg}.GetImage()
			if err != nil {
				glog.Error(err.Error())
				os.Exit(1)
			}
		}()

		go func() {
			defer wg.Done()
			image2, err = utils.ImagePrepper{img2Arg}.GetImage()
			if err != nil {
				glog.Error(err.Error())
				os.Exit(1)
			}
		}()

		diffTypes, err := differs.GetDiffers(diffArgs)
		if err != nil {
			glog.Error(err.Error())
			os.Exit(1)
		}
		wg.Wait()

		req := differs.DiffRequest{image1, image2, diffTypes}
		if diffs, err := req.GetDiff(); err == nil {
			// Outputs diff results in alphabetical order by differ name
			diffTypes := []string{}
			for name := range diffs {
				diffTypes = append(diffTypes, name)
			}
			sort.Strings(diffTypes)
			glog.Info("Retrieving diffs")
			diffResults := []utils.DiffResult{}
			for _, diffType := range diffTypes {
				diff := diffs[diffType]
				if json {
					diffResults = append(diffResults, diff.GetStruct())
				} else {
					err = diff.OutputText(diffType)
					if err != nil {
						glog.Error(err)
					}
				}
			}
			if json {
				err = utils.JSONify(diffResults)
				if err != nil {
					glog.Error(err)
				}
			}
			fmt.Println()
			glog.Info("Removing image file system directories from system")
			if !save {
				errMsg := remove(image1.FSPath, true)
				errMsg += remove(image2.FSPath, true)
				if errMsg != "" {
					glog.Error(errMsg)
				}
			}
		} else {
			glog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func getAllDiffers() []string {
	allDiffers := []string{}
	for name := range diffFlagMap {
		allDiffers = append(allDiffers, name)
	}
	return allDiffers
}

func validateArgs(args []string) (bool, error) {
	validArgNum, err := checkArgNum(args)
	if err != nil {
		return false, err
	} else if !validArgNum {
		return false, nil
	}
	validArgType, err := checkArgType(args)
	if err != nil {
		return false, err
	} else if !validArgType {
		return false, nil
	}
	return true, nil
}

func checkArgNum(args []string) (bool, error) {
	var errMessage string
	if len(args) < 2 {
		errMessage = "Too few arguments. Should have two images as arguments: [IMAGE1] [IMAGE2]."
		return false, errors.New(errMessage)
	} else if len(args) > 2 {
		errMessage = "Too many arguments. Should have two images as arguments: [IMAGE1] [IMAGE2]."
		return false, errors.New(errMessage)
	} else {
		return true, nil
	}
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
	if !checkImage(args[0]) {
		valid = false
		errMessage := fmt.Sprintf("Argument %s is not an image ID, URL, or tar\n", args[0])
		buffer.WriteString(errMessage)
	}
	if !checkImage(args[1]) {
		valid = false
		errMessage := fmt.Sprintf("Argument %s is not an image ID, URL, or tar\n", args[1])
		buffer.WriteString(errMessage)
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
	RootCmd.Flags().BoolVarP(&json, "json", "j", false, "JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).")
	RootCmd.Flags().BoolVarP(&eng, "eng", "e", false, "By default the docker calls are shelled out locally, set this flag to use the Docker Engine Client (version compatibility required).")
	RootCmd.Flags().BoolVarP(&pip, "pip", "p", false, "Set this flag to use the pip differ.")
	RootCmd.Flags().BoolVarP(&node, "node", "n", false, "Set this flag to use the node differ.")
	RootCmd.Flags().BoolVarP(&apt, "apt", "a", false, "Set this flag to use the apt differ.")
	RootCmd.Flags().BoolVarP(&file, "file", "f", false, "Set this flag to use the file differ.")
	RootCmd.Flags().BoolVarP(&history, "history", "d", false, "Set this flag to use the dockerfile history differ.")
	RootCmd.Flags().BoolVarP(&save, "save", "s", false, "Set this flag to save rather than remove the final image filesystems on exit.")
}
