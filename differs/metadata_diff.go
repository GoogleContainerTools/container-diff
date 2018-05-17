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

package differs

import (
	"fmt"
	"strings"

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
)

type MetadataAnalyzer struct {
}

type MetadataDiff struct {
	Adds []string
	Dels []string
}

func (a MetadataAnalyzer) Name() string {
	return "MetadataAnalyzer"
}

func (a MetadataAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff, err := getMetadataDiff(image1, image2)
	return &util.MetadataDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: "Metadata",
		Diff:     diff,
	}, err
}

func (a MetadataAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	analysis, err := getMetadataList(image)
	if err != nil {
		return &util.ListAnalyzeResult{}, err
	}
	return &util.ListAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: "Metadata",
		Analysis:    analysis,
	}, nil
}

func getMetadataDiff(image1, image2 pkgutil.Image) (MetadataDiff, error) {
	m1, err := getMetadataList(image1)
	if err != nil {
		return MetadataDiff{}, err
	}
	m2, err := getMetadataList(image2)
	if err != nil {
		return MetadataDiff{}, err
	}

	adds := util.GetAdditions(m1, m2)
	dels := util.GetDeletions(m1, m2)
	return MetadataDiff{adds, dels}, nil
}

func getMetadataList(image pkgutil.Image) ([]string, error) {
	configFile, err := image.Image.ConfigFile()
	if err != nil {
		return nil, err
	}
	c := configFile.Config

	return []string{
		fmt.Sprintf("Domainname: %s", c.Domainname),
		fmt.Sprintf("User: %s", c.User),
		fmt.Sprintf("AttachStdin: %t", c.AttachStdin),
		fmt.Sprintf("AttachStdout: %t", c.AttachStdout),
		fmt.Sprintf("AttachStderr: %t", c.AttachStderr),
		fmt.Sprintf("ExposedPorts: %v", pkgutil.SortMap(StructMapToStringMap(c.ExposedPorts))),
		fmt.Sprintf("Tty: %t", c.Tty),
		fmt.Sprintf("OpenStdin: %t", c.OpenStdin),
		fmt.Sprintf("StdinOnce: %t", c.StdinOnce),
		fmt.Sprintf("Env: %s", strings.Join(c.Env, ",")),
		fmt.Sprintf("Cmd: %s", strings.Join(c.Cmd, ",")),
		fmt.Sprintf("ArgsEscaped: %t", c.ArgsEscaped),
		fmt.Sprintf("Volumes: %v", pkgutil.SortMap(StructMapToStringMap(c.Volumes))),
		fmt.Sprintf("Workdir: %s", c.WorkingDir),
		fmt.Sprintf("Entrypoint: %s", strings.Join(c.Entrypoint, ",")),
		fmt.Sprintf("NetworkDisabled: %t", c.NetworkDisabled),
		fmt.Sprintf("MacAddress: %s", c.MacAddress),
		fmt.Sprintf("OnBuild: %s", strings.Join(c.OnBuild, ",")),
		fmt.Sprintf("Labels: %v", pkgutil.SortMap(c.Labels)),
		fmt.Sprintf("StopSignal: %s", c.StopSignal),
		fmt.Sprintf("Shell: %s", strings.Join(c.Shell, ",")),
	}, nil
}

// StructMapToStringMap converts map[string]struct{} to map[string]string knowing that the
// struct in the value is always empty
func StructMapToStringMap(m map[string]struct{}) map[string]string {
	newMap := make(map[string]string)
	for k := range m {
		newMap[k] = fmt.Sprintf("%s", m[k])
	}
	return newMap
}
