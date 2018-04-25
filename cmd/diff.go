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
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/GoogleContainerTools/container-diff/cmd/util/output"
	"github.com/GoogleContainerTools/container-diff/differs"
	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var filename string

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare two images: [image1] [image2]",
	Long:  `Compares two images using the specifed analyzers as indicated via flags (see documentation for available ones).`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := validateArgs(args, checkDiffArgNum, checkIfValidAnalyzer, checkFilenameFlag); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := diffImages(args[0], args[1], types); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	},
}

func checkDiffArgNum(args []string) error {
	if len(args) != 2 {
		return errors.New("'diff' requires two images as arguments: container-diff diff [image1] [image2]")
	}
	return nil
}

func checkFilenameFlag(_ []string) error {
	if filename == "" {
		return nil
	}
	for _, t := range types {
		if t == "file" {
			return nil
		}
	}
	return errors.New("Please include --types=file with the --filename flag")
}

func diffImages(image1Arg, image2Arg string, diffArgs []string) error {
	diffTypes, err := differs.GetAnalyzers(diffArgs)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	output.PrintToStdErr("Starting diff on images %s and %s, using differs: %s\n", image1Arg, image2Arg, diffArgs)

	imageMap := map[string]*pkgutil.Image{
		image1Arg: {},
		image2Arg: {},
	}
	// TODO: fix error handling here
	for imageArg := range imageMap {
		go func(imageName string, imageMap map[string]*pkgutil.Image) {
			defer wg.Done()
			image, err := getImageForName(imageName)
			imageMap[imageName] = &image
			if err != nil {
				logrus.Warningf("Diff may be inaccurate: %s", err)
			}
		}(imageArg, imageMap)
	}
	wg.Wait()

	if !save {
		defer pkgutil.CleanupImage(*imageMap[image1Arg])
		defer pkgutil.CleanupImage(*imageMap[image2Arg])
	}

	output.PrintToStdErr("Computing diffs\n")
	req := differs.DiffRequest{
		Image1:    *imageMap[image1Arg],
		Image2:    *imageMap[image2Arg],
		DiffTypes: diffTypes}
	diffs, err := req.GetDiff()
	if err != nil {
		return fmt.Errorf("Could not retrieve diff: %s", err)
	}
	outputResults(diffs)

	if filename != "" {
		output.PrintToStdErr("Computing filename diffs\n")
		err := diffFile(imageMap[image1Arg], imageMap[image2Arg])
		if err != nil {
			return err
		}
	}

	if save {
		logrus.Infof("Images were saved at %s and %s", imageMap[image1Arg].FSPath,
			imageMap[image2Arg].FSPath)
	}
	return nil
}

func diffFile(image1, image2 *pkgutil.Image) error {
	diff, err := util.DiffFile(image1, image2, filename)
	if err != nil {
		return err
	}
	util.TemplateOutput(diff, "FilenameDiff")
	return nil
}

func init() {
	diffCmd.Flags().StringVarP(&filename, "filename", "f", "", "Set this flag to the path of a file in both containers to view the diff of the file. Must be used with --types=file flag.")
	RootCmd.AddCommand(diffCmd)
	addSharedFlags(diffCmd)
	output.AddFlags(diffCmd)
}
