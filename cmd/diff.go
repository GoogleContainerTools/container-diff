package cmd

import (
	"errors"
	"os"
	"sync"

	"github.com/GoogleCloudPlatform/container-diff/differs"
	"github.com/GoogleCloudPlatform/container-diff/utils"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare two images: [image1] [image2]",
	Long:  `Compares two images using the specifed analyzers as indicated via flags (see documentation for available ones).`,
	Run: func(cmd *cobra.Command, args []string) {
		if validArgs, err := validateArgs(args, checkDiffArgNum, checkArgType); !validArgs {
			glog.Error(err.Error())
			os.Exit(1)
		}

		utils.SetDockerEngine(eng)

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

		if err := diffImages(args[0], args[1], analyzeArgs); err != nil {
			glog.Error(err)
			os.Exit(1)
		}
	},
}

func checkDiffArgNum(args []string) (bool, error) {
	var errMessage string
	if len(args) != 2 {
		errMessage = "'diff' requires two images as arguments: container diff [image1] [image2]"
		return false, errors.New(errMessage)
	}
	return true, nil
}

func diffImages(image1Arg, image2Arg string, diffArgs []string) error {
	var wg sync.WaitGroup
	wg.Add(2)

	glog.Infof("Starting diff on images %s and %s, using differs: %s", image1Arg, image2Arg, diffArgs)

	var image1, image2 utils.Image
	var err error
	go func() {
		defer wg.Done()
		image1, err = utils.ImagePrepper{image1Arg}.GetImage()
		if err != nil {
			glog.Error(err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		image2, err = utils.ImagePrepper{image2Arg}.GetImage()
		if err != nil {
			glog.Error(err.Error())
		}
	}()
	wg.Wait()
	if err != nil {
		cleanupImage(image1)
		cleanupImage(image2)
		return errors.New("Could not perform image diff")
	}

	diffTypes, err := differs.GetAnalyzers(diffArgs)
	if err != nil {
		glog.Error(err.Error())
		cleanupImage(image1)
		cleanupImage(image2)
		return errors.New("Could not perform image diff")
	}

	req := differs.DiffRequest{image1, image2, diffTypes}
	if diffs, err := req.GetDiff(); err == nil {
		glog.Info("Retrieving diffs")
		outputResults(diffs)
		if !save {
			cleanupImage(image1)
			cleanupImage(image2)

		} else {
			dir, _ := os.Getwd()
			glog.Infof("Images were saved at %s as %s and %s", dir, image1.FSPath, image2.FSPath)
		}
	} else {
		glog.Error(err.Error())
		cleanupImage(image1)
		cleanupImage(image2)
		return errors.New("Could not perform image diff")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(diffCmd)
	addSharedFlags(diffCmd)
}
