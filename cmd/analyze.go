/*
Copyright 2018 Google, Inc. All rights reserved.

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
	"fmt"
	"os"

	"github.com/GoogleContainerTools/container-diff/cmd/util/output"
	"github.com/GoogleContainerTools/container-diff/differs"
	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze image",
	Short: "Analyzes an image: container-diff analyze image",
	Long: `Analyzes an image using the specifed analyzers as indicated via --type flag(s).

For details on how to specify images, run: container-diff help`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := validateArgs(args, checkAnalyzeArgNum, checkIfValidAnalyzer); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := analyzeImage(args[0], types); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	},
}

func checkAnalyzeArgNum(args []string) error {
	if len(args) != 1 {
		return errors.New("'analyze' requires one image as an argument: container-diff analyze [image]")
	}
	return nil
}

func analyzeImage(imageName string, analyzerArgs []string) error {
	analyzeTypes, err := differs.GetAnalyzers(analyzerArgs)
	if err != nil {
		return errors.Wrap(err, "getting analyzers")
	}

	image, err := getImage(imageName)
	if err != nil {
		return errors.Wrapf(err, "error retrieving image %s", imageName)
	}

	if noCache && !save {
		defer pkgutil.CleanupImage(image)
	}
	if err != nil {
		return fmt.Errorf("error processing image: %s", err)
	}

	req := differs.SingleRequest{
		Image:        image,
		AnalyzeTypes: analyzeTypes}
	analyses, err := req.GetAnalysis()
	if err != nil {
		return fmt.Errorf("error performing image analysis: %s", err)
	}

	logrus.Info("retrieving analyses")
	outputResults(analyses)

	if noCache && save {
		logrus.Infof("image was saved at %s", image.FSPath)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(analyzeCmd)
	addSharedFlags(analyzeCmd)
	output.AddFlags(analyzeCmd)
}
