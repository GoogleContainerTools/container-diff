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
	"context"
	goflag "flag"
	"fmt"
	"sort"
	"strings"

	"github.com/GoogleCloudPlatform/container-diff/differs"
	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/GoogleCloudPlatform/container-diff/util"
	"github.com/docker/docker/client"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var json bool
var save bool
var types string

type validatefxn func(args []string) error

var RootCmd = &cobra.Command{
	Use:   "container-diff",
	Short: "container-diff is a tool for analyzing and comparing container images",
	Long:  `container-diff is a CLI tool for analyzing and comparing container images.`,
}

func NewClient() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("Error getting docker client: %s", err)
	}
	cli.NegotiateAPIVersion(context.Background())

	return cli, nil
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
			err := result.OutputText(analyzerType)
			if err != nil {
				glog.Error(err)
			}
		}
	}
	if json {
		err := util.JSONify(results)
		if err != nil {
			glog.Error(err)
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

func checkImage(arg string) bool {
	if !pkgutil.CheckImageID(arg) && !pkgutil.CheckImageURL(arg) && !pkgutil.CheckTar(arg) {
		return false
	}
	return true
}

func checkArgType(args []string) error {
	for _, arg := range args {
		if !checkImage(arg) {
			errMessage := fmt.Sprintf("Argument %s is not an image ID, URL, or tar\n", args[0])
			glog.Errorf(errMessage)
			return errors.New(errMessage)
		}
	}
	return nil
}

func checkIfValidAnalyzer(flagtypes string) error {
	analyzers := strings.Split(flagtypes, ",")
	for _, name := range analyzers {
		if _, exists := differs.Analyzers[name]; !exists {
			errMessage := fmt.Sprintf("Argument %s is not an image ID, URL, or tar\n", name)
			glog.Errorf(errMessage)
			return errors.New(errMessage)
		}
	}
	return nil
}

func init() {
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func addSharedFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&json, "json", "j", false, "JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).")
	cmd.Flags().StringVarP(&types, "types", "t", "", "This flag sets the list of analyzer types to use.  It expects a comma separated list of supported analyzers.")
	cmd.Flags().BoolVarP(&save, "save", "s", false, "Set this flag to save rather than remove the final image filesystems on exit.")
	cmd.Flags().BoolVarP(&util.SortSize, "order", "o", false, "Set this flag to sort any file/package results by descending size. Otherwise, they will be sorted by name.")
}
