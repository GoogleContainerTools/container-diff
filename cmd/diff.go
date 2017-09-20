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
	"sync"

	"github.com/GoogleCloudPlatform/container-diff/differs"
	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
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
		if err := diffImages(args[0], args[1], strings.Split(types, ",")); err != nil {
			glog.Error(err)
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

func diffImages(image1Arg, image2Arg string, diffArgs []string) error {
	diffTypes, err := differs.GetAnalyzers(diffArgs)
	if err != nil {
		return fmt.Errorf("err msg: %s", err.Error())
	}

	cli, err := NewClient()
	if err != nil {
		return fmt.Errorf("err msg: %s", err)
	}
	defer cli.Close()
	var wg sync.WaitGroup
	wg.Add(2)

	glog.Infof("Starting diff on images %s and %s, using differs: %s", image1Arg, image2Arg, diffArgs)

	imageMap := map[string]*pkgutil.Image{
		image1Arg: {},
		image2Arg: {},
	}
	for imageArg := range imageMap {
		go func(imageName string, imageMap map[string]*pkgutil.Image) {
			defer wg.Done()

			prefixedName := processImageName(imageName)
			ip := pkgutil.ImagePrepper{
				Source: prefixedName,
				Client: cli,
			}
			image, err := ip.GetImage()
			imageMap[imageName] = &image
			if err != nil {
				glog.Errorf("Diff may be inaccurate: %s", err.Error())
			}
		}(imageArg, imageMap)
	}
	wg.Wait()

	if !save {
		defer pkgutil.CleanupImage(*imageMap[image1Arg])
		defer pkgutil.CleanupImage(*imageMap[image2Arg])
	}

	req := differs.DiffRequest{*imageMap[image1Arg], *imageMap[image2Arg], diffTypes}
	diffs, err := req.GetDiff()
	if err != nil {
		return fmt.Errorf("err msg: %s", err.Error())
	}
	glog.Info("Retrieving diffs")
	outputResults(diffs)

	if save {
		glog.Infof("Images were saved at %s and %s", imageMap[image1Arg].FSPath,
			imageMap[image2Arg].FSPath)

	}
	return nil
}

func processImageName(imageName string) string {
	if pkgutil.IsTar(imageName) || strings.HasPrefix(imageName, "daemon://") {
		return imageName
	}
	// not a tar and not explicitly local, so force remote
	return "remote://" + imageName
}

func init() {
	RootCmd.AddCommand(diffCmd)
	addSharedFlags(diffCmd)
}
