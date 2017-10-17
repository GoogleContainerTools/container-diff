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
	"github.com/GoogleCloudPlatform/container-diff/differs"
	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/GoogleCloudPlatform/container-diff/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"sync"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare two images: [image1] [image2]",
	Long:  `Compares two images using the specifed analyzers as indicated via flags (see documentation for available ones).`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := validateArgs(args, checkDiffArgNum); err != nil {
			return err
		}
		if err := checkIfValidAnalyzer(types); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		typesFlagSet := checkIfTypesFlagSet(cmd)
		if err := diffImages(args[0], args[1], strings.Split(types, ","), typesFlagSet); err != nil {
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

func diffImages(image1Arg, image2Arg string, diffArgs []string, typesFlagSet bool) error {
	diffTypes, err := differs.GetAnalyzers(diffArgs)
	if err != nil {
		return err
	}

	cli, err := pkgutil.NewClient()
	if err != nil {
		return err
	}
	defer cli.Close()
	var wg sync.WaitGroup
	wg.Add(2)

	if typesFlagSet || filename == "" {
		fmt.Fprintf(os.Stderr, "Starting diff on images %s and %s, using differs: %s\n", image1Arg, image2Arg, diffArgs)
	}

	imageMap := map[string]*pkgutil.Image{
		image1Arg: {},
		image2Arg: {},
	}
	// TODO: fix error handling here
	for imageArg := range imageMap {
		go func(imageName string, imageMap map[string]*pkgutil.Image) {
			defer wg.Done()

			prepper, err := getPrepperForImage(imageName)
			if err != nil {
				logrus.Error(err)
				return
			}
			image, err := prepper.GetImage()
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

	if filename != "" {
		fmt.Fprintln(os.Stderr, "Computing filename diffs")
		err := diffFile(imageMap, image1Arg, image2Arg)
		if err != nil {
			return err
		}
		if !typesFlagSet {
			return nil
		}
	}

	fmt.Fprintln(os.Stderr, "Computing diffs")
	req := differs.DiffRequest{*imageMap[image1Arg], *imageMap[image2Arg], diffTypes}
	diffs, err := req.GetDiff()
	if err != nil {
		return fmt.Errorf("Could not retrieve diff: %s", err)
	}
	outputResults(diffs)

	if save {
		logrus.Infof("Images were saved at %s and %s", imageMap[image1Arg].FSPath,
			imageMap[image2Arg].FSPath)
	}
	return nil
}

func diffFile(imageMap map[string]*pkgutil.Image, image1Arg, image2Arg string) error {

	image1FilePath := imageMap[image1Arg].FSPath + filename
	image2FilePath := imageMap[image2Arg].FSPath + filename

	diff, err := util.DiffFile(image1FilePath, image2FilePath, image1Arg, image2Arg)
	if err != nil {
		return err
	}
	diff.Filename = filename
	util.TemplateOutput(diff, "FileNameDiff")
	return nil
}

func init() {
	RootCmd.AddCommand(diffCmd)
	addSharedFlags(diffCmd)
}
