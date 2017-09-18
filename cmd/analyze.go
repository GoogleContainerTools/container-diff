/*
Copyright 2017 Google, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/container-diff/differs"
	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyzes an image: [image]",
	Long:  `Analyzes an image using the specifed analyzers as indicated via flags (see documentation for available ones).`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := validateArgs(args, checkAnalyzeArgNum, checkArgType); err != nil {
			return errors.New(err.Error())
		}
		if err := checkIfValidAnalyzer(types); err != nil {
			return errors.New(err.Error())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := analyzeImage(args[0], strings.Split(types, ",")); err != nil {
			glog.Error(err)
			os.Exit(1)
		}
	},
}

func checkAnalyzeArgNum(args []string) error {
	if len(args) != 1 {
		return errors.New("'analyze' requires one image as an argument: container analyze [image]")
	}
	return nil
}

func analyzeImage(imageName string, analyzerArgs []string) error {
	cli, err := NewClient()
	if err != nil {
		return fmt.Errorf("Error getting docker client for differ: %s", err)
	}
	defer cli.Close()

	var prepper string
	if strings.HasPrefix(imageName, "daemon://") {
		// force local daemon if we have the corresponding prefix
		prepper = pkgutil.LOCAL
		imageName = strings.Replace(imageName, "daemon://", "", -1)
	} else if pkgutil.IsTar(imageName) {
		prepper = pkgutil.TAR
	} else {
		prepper = pkgutil.REMOTE
	}

	ip := pkgutil.ImagePrepper{
		Source:  imageName,
		Client:  cli,
		Prepper: prepper,
	}
	image, err := ip.GetImage()

	if !save {
		defer pkgutil.CleanupImage(image)
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
		glog.Infof("Image was saved at %s", image.FSPath)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(analyzeCmd)
	addSharedFlags(analyzeCmd)
}
