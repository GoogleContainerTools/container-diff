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
	"strings"
	"sync"

	"github.com/GoogleContainerTools/container-diff/cmd/util/output"
	"github.com/GoogleContainerTools/container-diff/differs"
	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
	"github.com/pkg/errors"
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

func processImage(imageName string, imageMap map[string]*pkgutil.Image, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	image, err := getImageForName(imageName)
	if image.Image == nil {
		errChan <- fmt.Errorf("error retrieving image %s: %s", imageName, err.Error())
		return
	}
	if err != nil {
		logrus.Warningf("diff may be inaccurate: %s", err)
	}
	imageMap[imageName] = &image
}

func diffImages(image1Arg, image2Arg string, diffArgs []string) error {
	diffTypes, err := differs.GetAnalyzers(diffArgs)
	if err != nil {
		return errors.Wrap(err, "getting analyzers")
	}

	var wg sync.WaitGroup
	wg.Add(2)

	logrus.Infof("starting diff on images %s and %s, using differs: %s\n", image1Arg, image2Arg, diffArgs)

	imageMap := map[string]*pkgutil.Image{}
	errChan := make(chan error, 2)

	go processImage(image1Arg, imageMap, &wg, errChan)
	go processImage(image2Arg, imageMap, &wg, errChan)

	wg.Wait()
	close(errChan)

	errs := []string{}
	for {
		err, ok := <-errChan
		if !ok {
			break
		}
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	img1, ok := imageMap[image1Arg]
	if !ok {
		return fmt.Errorf("cannot find image %s", image1Arg)
	}
	img2, ok := imageMap[image2Arg]
	if !ok {
		return fmt.Errorf("cannot find image %s", image2Arg)
	}

	if noCache && !save {
		defer pkgutil.CleanupImage(*imageMap[image1Arg])
		defer pkgutil.CleanupImage(*imageMap[image2Arg])
	}

	logrus.Info("computing diffs")
	req := differs.DiffRequest{
		Image1:    *img1,
		Image2:    *img2,
		DiffTypes: diffTypes}
	diffs, err := req.GetDiff()
	if err != nil {
		return fmt.Errorf("Could not retrieve diff: %s", err)
	}
	outputResults(diffs)

	if filename != "" {
		logrus.Info("computing filename diffs")
		err := diffFile(imageMap[image1Arg], imageMap[image2Arg])
		if err != nil {
			return err
		}
	}

	if noCache && save {
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
