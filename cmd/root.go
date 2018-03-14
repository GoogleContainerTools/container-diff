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
	goflag "flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GoogleCloudPlatform/container-diff/differs"
	"github.com/GoogleCloudPlatform/container-diff/pkg/cache"
	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/GoogleCloudPlatform/container-diff/util"
	"github.com/containers/image/docker"
	"github.com/docker/docker/client"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var json bool
var save bool
var types diffTypes
var noCache bool

var LogLevel string
var format string

type validatefxn func(args []string) error

const (
	DaemonPrefix = "daemon://"
	RemotePrefix = "remote://"
)

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

	results := make([]interface{}, len(resultMap))
	for i, analyzerType := range sortedTypes {
		result := resultMap[analyzerType]
		if json {
			results[i] = result.OutputStruct()
		} else {
			err := result.OutputText(analyzerType, format)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
	if json {
		err := util.JSONify(results)
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
		types = []string{"apt"}
	}
	for _, name := range types {
		if _, exists := differs.Analyzers[name]; !exists {
			return fmt.Errorf("Argument %s is not a valid analyzer", name)
		}
	}
	return nil
}

func getPrepperForImage(image string) (pkgutil.Prepper, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	if pkgutil.IsTar(image) {
		return &pkgutil.TarPrepper{
			Source: image,
			Client: cli,
		}, nil
	}

	if strings.HasPrefix(image, DaemonPrefix) {
		// remove the DaemonPrefix
		image = strings.Replace(image, DaemonPrefix, "", -1)

		return &pkgutil.DaemonPrepper{
			Source: image,
			Client: cli,
		}, nil
	}

	// either has remote prefix or has no prefix, in which case we force remote
	ref, err := docker.ParseReference("//" + image)
	if err != nil {
		return nil, err
	}
	src, err := ref.NewImageSource(nil)
	if err != nil {
		return nil, err
	}

	if !noCache {
		cacheDir, err := cacheDir()
		if err != nil {
			return nil, err
		}

		src, err = cache.NewFileCache(cacheDir, ref)
		if err != nil {
			return nil, err
		}
	}

	return &pkgutil.CloudPrepper{
		Source:      strings.Replace(image, RemotePrefix, "", -1),
		Client:      cli,
		ImageSource: src,
	}, nil
}

func cacheDir() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	rootDir := filepath.Join(dir, ".container-diff")
	return filepath.Join(rootDir, "cache"), nil
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
	cmd.Flags().BoolVarP(&json, "json", "j", false, "JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).")
	cmd.Flags().VarP(&types, "type", "t", "This flag sets the list of analyzer types to use. Set it repeatedly to use multiple analyzers.")
	cmd.Flags().BoolVarP(&save, "save", "s", false, "Set this flag to save rather than remove the final image filesystems on exit.")
	cmd.Flags().BoolVarP(&util.SortSize, "order", "o", false, "Set this flag to sort any file/package results by descending size. Otherwise, they will be sorted by name.")
	cmd.Flags().BoolVarP(&noCache, "no-cache", "n", false, "Set this to force retrieval of layers on each run.")
}
