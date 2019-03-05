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
	goflag "flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GoogleContainerTools/container-diff/differs"
	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var json bool

var save bool
var types diffTypes
var noCache bool

var outputFile string
var forceWrite bool
var cacheDir string
var LogLevel string
var format string

const containerDiffEnvCacheDir = "CONTAINER_DIFF_CACHEDIR"

type validatefxn func(args []string) error

var RootCmd = &cobra.Command{
	Use:   "container-diff",
	Short: "container-diff is a tool for analyzing and comparing container images",
	Long: `container-diff is a CLI tool for analyzing and comparing container images.

Images can be specified from either a local Docker daemon, or from a remote registry.
To specify a local image, prefix the image ID with 'daemon://', e.g. 'daemon://gcr.io/foo/bar'.
To specify a remote image, prefix the image ID with 'remote://', e.g. 'remote://gcr.io/foo/bar'.
If no prefix is specified, the local daemon will be checked first.

Tarballs can also be specified by simply providing the path to the .tar, .tar.gz, or .tgz file.`,
	PersistentPreRun: func(c *cobra.Command, s []string) {
		ll, err := logrus.ParseLevel(LogLevel)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		logrus.SetLevel(ll)
	},
}

func outputResults(resultMap map[string]util.Result) {
	// Outputs diff/analysis results in alphabetical order by analyzer name
	sortedTypes := []string{}
	for analyzerType := range resultMap {
		sortedTypes = append(sortedTypes, analyzerType)
	}
	sort.Strings(sortedTypes)

	// Get the writer
	writer, err := getWriter(outputFile)
	if err != nil {
		errors.Wrap(err, "getting writer for output file")
	}

	results := make([]interface{}, len(resultMap))
	for i, analyzerType := range sortedTypes {
		result := resultMap[analyzerType]
		if json {
			results[i] = result.OutputStruct()
		} else {
			err := result.OutputText(writer, analyzerType, format)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
	if json {
		err := util.JSONify(writer, results)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func validateArgs(args []string, validatefxns ...validatefxn) error {
	for _, validatefxn := range validatefxns {
		if err := validatefxn(args); err != nil {
			return err
		}
	}
	return nil
}

func checkIfValidAnalyzer(_ []string) error {
	if len(types) == 0 {
		types = []string{"size"}
	}
	for _, name := range types {
		if _, exists := differs.Analyzers[name]; !exists {
			return fmt.Errorf("Argument %s is not a valid analyzer", name)
		}
	}
	return nil
}

func includeLayers() bool {
	for _, t := range types {
		for _, a := range differs.LayerAnalyzers {
			if t == a {
				return true
			}
		}
	}
	return false
}

func getImage(imageName string) (pkgutil.Image, error) {
	var cachePath string
	var err error
	if !noCache {
		cachePath, err = getCacheDir(imageName)
		if err != nil {
			return pkgutil.Image{}, err
		}
	}
	return pkgutil.GetImage(imageName, includeLayers(), cachePath)
}

func getCacheDir(imageName string) (string, error) {
	// First preference for cache is set at command line
	if cacheDir == "" {
		// second preference is environment
		cacheDir = os.Getenv(containerDiffEnvCacheDir)
	}

	// Third preference (default) is set at $HOME
	if cacheDir == "" {
		dir, err := homedir.Dir()
		if err != nil {
			return "", errors.Wrap(err, "retrieving home dir")
		} else {
			cacheDir = dir
		}
	}
	rootDir := filepath.Join(cacheDir, ".container-diff", "cache")
	imageName = strings.Replace(imageName, string(os.PathSeparator), "", -1)
	return filepath.Join(rootDir, pkgutil.CleanFilePath(imageName)), nil
}

func getWriter(outputFile string) (io.Writer, error) {
	var err error
	var outWriter io.Writer
	// If the user specifies an output file, ensure exists
	if outputFile != "" {
		// Don't overwrite a file that exists, unless given --force
		if _, err := os.Stat(outputFile); !os.IsNotExist(err) && !forceWrite {
			errors.Wrap(err, "file exist, will not overwrite.")
		}
		// Otherwise, output file is an io.writer
		outWriter, err = os.Create(outputFile)
	}
	// If still doesn't exist, return stdout as the io.Writer
	if outputFile == "" {
		outWriter = os.Stdout
	}
	return outWriter, err
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&LogLevel, "verbosity", "v", "warning", "This flag controls the verbosity of container-diff.")
	RootCmd.PersistentFlags().StringVarP(&format, "format", "", "", "Format to output diff in.")
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

// Define a type named "diffSlice" as a slice of strings
type diffTypes []string

// Now, for our new type, implement the two methods of
// the flag.Value interface...
// The first method is String() string
func (d *diffTypes) String() string {
	return strings.Join(*d, ",")
}

// The second method is Set(value string) error
func (d *diffTypes) Set(value string) error {
	// Dedupe repeated elements.
	for _, t := range *d {
		if t == value {
			return nil
		}
	}
	*d = append(*d, value)
	return nil
}

func (d *diffTypes) Type() string {
	return "Diff Types"
}

func addSharedFlags(cmd *cobra.Command) {
	sortedTypes := []string{}
	for analyzerType := range differs.Analyzers {
		sortedTypes = append(sortedTypes, analyzerType)
	}
	sort.Strings(sortedTypes)
	supportedTypes := strings.Join(sortedTypes, ", ")

	cmd.Flags().BoolVarP(&json, "json", "j", false, "JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).")
	cmd.Flags().VarP(&types, "type", "t",
		fmt.Sprintf("This flag sets the list of analyzer types to use.\n"+
			"Set it repeatedly to use multiple analyzers.\n"+
			"Supported types: %s.",
			supportedTypes))
	cmd.Flags().BoolVarP(&save, "save", "s", false, "Set this flag to save rather than remove the final image filesystems on exit.")
	cmd.Flags().BoolVarP(&util.SortSize, "order", "o", false, "Set this flag to sort any file/package results by descending size. Otherwise, they will be sorted by name.")
	cmd.Flags().BoolVarP(&noCache, "no-cache", "n", false, "Set this to force retrieval of image filesystem on each run.")
	cmd.Flags().StringVarP(&cacheDir, "cache-dir", "c", "", "cache directory base to create .container-diff (default is $HOME).")
	cmd.Flags().StringVarP(&outputFile, "output", "w", "", "output file to write to (default writes to the screen).")
	cmd.Flags().BoolVar(&forceWrite, "force", false, "force overwrite output file, if exists already.")
}
