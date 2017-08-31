package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/container-diff/differs"
	"github.com/GoogleCloudPlatform/container-diff/utils"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyzes an image: [image]",
	Long:  `Analyzes an image using the specifed analyzers as indicated via flags (see documentation for available ones).`,
	Run: func(cmd *cobra.Command, args []string) {
		if validArgs, err := validateArgs(args, checkAnalyzeArgNum, checkArgType); !validArgs {
			glog.Error(err.Error())
			os.Exit(1)
		}
		analyzeArgs := []string{}
		allAnalyzers := getAllAnalyzers()
		for _, name := range allAnalyzers {
			if *analyzeFlagMap[name] == true {
				analyzeArgs = append(analyzeArgs, name)
			}
		}

		// If no analyzers are specified, perform them all as the default
		if len(analyzeArgs) == 0 {
			analyzeArgs = allAnalyzers
		}

		if err := analyzeImage(args[0], analyzeArgs); err != nil {
			glog.Error(err)
			os.Exit(1)
		}
	},
}

func checkAnalyzeArgNum(args []string) (bool, error) {
	var errMessage string
	if len(args) != 1 {
		errMessage = "'analyze' requires one image as an argument: container analyze [image]"
		return false, errors.New(errMessage)
	}
	return true, nil
}

func analyzeImage(imageArg string, analyzerArgs []string) error {
	cli, err := NewClient()
	if err != nil {
		return fmt.Errorf("Error getting docker client for differ: %s", err)
	}
	defer cli.Close()
	ip := utils.ImagePrepper{
		Source: imageArg,
		Client: cli,
	}
	image, err := ip.GetImage()

	if !save {
		defer cleanupImage(image)
	}
	if err != nil {
		glog.Error(err.Error())
		return errors.New("Could not perform image analysis")
	}
	analyzeTypes, err := differs.GetAnalyzers(analyzerArgs)
	if err != nil {
		glog.Error(err.Error())
		return errors.New("Could not perform image analysis")
	}

	req := differs.SingleRequest{image, analyzeTypes}
	analyses, err := req.GetAnalysis()
	if err != nil {
		glog.Error(err.Error())
		return errors.New("Could not perform image analysis")
	}

	glog.Info("Retrieving analyses")
	outputResults(analyses)

	if save {
		dir, _ := os.Getwd()
		glog.Infof("Image was saved at %s as %s", dir, image.FSPath)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(analyzeCmd)
	addSharedFlags(analyzeCmd)
}
