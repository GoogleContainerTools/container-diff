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
var ci bool

var diffCmd = &cobra.Command{
	Use:   "diff image1 image2",
	Short: "Compare two images: container-diff diff image1 image2",
	Long: `Compares two images using the specifed analyzers as indicated via --type flag(s).

For details on how to specify images, run: container-diff help`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := validateArgs(args, checkDiffArgNum, checkDiffImages, checkIfValidAnalyzer, checkFilenameFlag); err != nil {
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

func checkDiffImages(args []string) error {
	if args[0] == args[1] {
		return errors.New("'diff' requires two different images")
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
	return errors.New("please include --types=file with the --filename flag")
}

// processImage is a concurrency-friendly wrapper around getImageForName
func processImage(imageName string, errChan chan<- error) *pkgutil.Image {
	image, err := getImage(imageName)
	if err != nil {
		errChan <- fmt.Errorf("error retrieving image %s: %s", imageName, err)
	}
	return &image
}

func processImages(image1Arg, image2Arg string) (*pkgutil.Image, *pkgutil.Image, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	var image1, image2 *pkgutil.Image
	errChan := make(chan error, 2)

	go func() {
		defer wg.Done()
		image1 = processImage(image1Arg, errChan)
	}()
	go func() {
		defer wg.Done()
		image2 = processImage(image2Arg, errChan)
	}()

	wg.Wait()
	close(errChan)

	if noCache && !save {
		defer pkgutil.CleanupImage(*image1)
		defer pkgutil.CleanupImage(*image2)
	}

	err := readErrorsFromChannel(errChan)

	return image1, image2, err
}

// collects errors from a channel and combines them
// assumes channel has already been closed
func readErrorsFromChannel(c chan error) error {
	errs := []string{}
	for err := range c {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func diffImages(image1Arg, image2Arg string, diffArgs []string) error {
	diffTypes, err := differs.GetAnalyzers(diffArgs)
	if err != nil {
		return errors.Wrap(err, "getting analyzers")
	}

	logrus.Infof("starting diff on images %s and %s, using differs: %s\n", image1Arg, image2Arg, diffArgs)
	image1, image2, err := processImages(image1Arg, image2Arg)
	if err != nil {
		return errors.Wrap(err, "Processing input images")
	}

	logrus.Info("computing diffs")
	req := differs.DiffRequest{
		Image1:    *image1,
		Image2:    *image2,
		DiffTypes: diffTypes}
	diffs, err := req.GetDiff()
	if err != nil {
		return fmt.Errorf("could not retrieve diff: %s", err)
	}
	outputResults(diffs)

	if filename != "" {
		logrus.Info("computing filename diffs")
		err := diffFile(image1, image2)
		if err != nil {
			return err
		}
	}

	if noCache && save {
		logrus.Infof("images were saved at %s and %s", image1.FSPath,
			image2.FSPath)
	}

	if ci {
		counter := getDiffCount(diffs)
		if counter > 0 {
			os.Exit(counter)
		}
	}

	return nil
}

func getDiffCount(diffs map[string]util.Result) int {
	var counter = 0
	for i, result := range diffs {
		switch i {
		case "HistoryAnalyzer":
			diff := result.(*util.HistDiffResult).Diff.(differs.HistDiff)
			counter += len(diff.Adds) + len(diff.Dels)
			break
		case "MetadataAnalyzer":
			diff := result.(*util.MetadataDiffResult).Diff.(differs.MetadataDiff)
			counter += len(diff.Adds) + len(diff.Dels)
			break
		case "RPMAnalyzer", "AptAnalyzer":
			diff := result.(*util.SingleVersionPackageDiffResult).Diff.(util.PackageDiff)
			counter += len(diff.Packages1) + len(diff.Packages2) + len(diff.InfoDiff)
			break
		case "FileAnalyzer":
			diff := result.(*util.DirDiffResult).Diff.(util.DirDiff)
			counter += len(diff.Adds) + len(diff.Dels) + len(diff.Mods)
			break
		}
	}
	return counter
}

func diffFile(image1, image2 *pkgutil.Image) error {
	diff, err := util.DiffFile(image1, image2, filename)
	if err != nil {
		return err
	}
	writer, err := getWriter(outputFile)
	if err != nil {
		return err
	}
	util.TemplateOutput(writer, diff, "FilenameDiff")
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func init() {
	diffCmd.Flags().StringVarP(&filename, "filename", "f", "", "Set this flag to the path of a file in both containers to view the diff of the file. Must be used with --types=file flag.")
	diffCmd.Flags().BoolVar(&ci, "ci", false, "If set will exit the with error code when differences where found.")
	RootCmd.AddCommand(diffCmd)
	addSharedFlags(diffCmd)
	output.AddFlags(diffCmd)
}
