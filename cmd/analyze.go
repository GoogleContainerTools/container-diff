package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
		if err := validateArgs(args, checkAnalyzeArgNum, checkArgType); err != nil {
			glog.Error(err.Error())
			os.Exit(1)
		}
		if err := checkIfValidAnalyzer(types); err != nil {
			glog.Error(err)
			os.Exit(1)
		}
		if err := analyzeImage(args[0], strings.Split(types, ",")); err != nil {
			glog.Error(err)
			os.Exit(1)
		}
	},
}

func checkAnalyzeArgNum(args []string) error {
	var errMessage string
	if len(args) != 1 {
		errMessage = "'analyze' requires one image as an argument: container analyze [image]"
		glog.Errorf(errMessage)
		return errors.New(errMessage)
	}
	return nil
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
	if err != nil {
		glog.Error(err.Error())
		cleanupImage(image)
		return errors.New("Could not perform image analysis")
	}
	analyzeTypes, err := differs.GetAnalyzers(analyzerArgs)
	if err != nil {
		glog.Error(err.Error())
		cleanupImage(image)
		return errors.New("Could not perform image analysis")
	}

	req := differs.SingleRequest{image, analyzeTypes}
	if analyses, err := req.GetAnalysis(); err == nil {
		glog.Info("Retrieving analyses")
		outputResults(analyses)
		if !save {
			cleanupImage(image)
		} else {
			dir, _ := os.Getwd()
			glog.Infof("Image was saved at %s as %s", dir, image.FSPath)
		}
	} else {
		glog.Error(err.Error())
		cleanupImage(image)
		return errors.New("Could not perform image analysis")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(analyzeCmd)
	addSharedFlags(analyzeCmd)
}
