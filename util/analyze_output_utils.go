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

package util

import (
	"errors"
	"fmt"

	"github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/sirupsen/logrus"
)

type Result interface {
	OutputStruct() interface{}
	OutputText(resultType string, format string) error
}

type AnalyzeResult struct {
	Image       string
	AnalyzeType string
	Analysis    interface{}
}

type ListAnalyzeResult AnalyzeResult

func (r ListAnalyzeResult) OutputStruct() interface{} {
	return r
}

func (r ListAnalyzeResult) OutputText(resultType string, format string) error {
	analysis, valid := r.Analysis.([]string)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type []string")
		return fmt.Errorf("Could not output %s analysis result", r.AnalyzeType)
	}
	r.Analysis = analysis
	return TemplateOutputFromFormat(r, "ListAnalyze", format)

}

type MultiVersionPackageAnalyzeResult AnalyzeResult

func (r MultiVersionPackageAnalyzeResult) OutputStruct() interface{} {
	analysis, valid := r.Analysis.(map[string]map[string]PackageInfo)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type map[string]map[string]PackageInfo")
		return fmt.Errorf("Could not output %s analysis result", r.AnalyzeType)
	}
	analysisOutput := getMultiVersionPackageOutput(analysis)
	output := struct {
		Image       string
		AnalyzeType string
		Analysis    []PackageOutput
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    analysisOutput,
	}
	return output
}

func (r MultiVersionPackageAnalyzeResult) OutputText(resultType string, format string) error {
	analysis, valid := r.Analysis.(map[string]map[string]PackageInfo)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type map[string]map[string]PackageInfo")
		return fmt.Errorf("Could not output %s analysis result", r.AnalyzeType)
	}
	analysisOutput := getMultiVersionPackageOutput(analysis)

	strAnalysis := stringifyPackages(analysisOutput)
	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    []StrPackageOutput
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    strAnalysis,
	}
	return TemplateOutputFromFormat(strResult, "MultiVersionPackageAnalyze", format)
}

type SingleVersionPackageAnalyzeResult AnalyzeResult

func (r SingleVersionPackageAnalyzeResult) OutputStruct() interface{} {
	analysis, valid := r.Analysis.(map[string]PackageInfo)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type map[string]PackageInfo")
		return fmt.Errorf("Could not output %s analysis result", r.AnalyzeType)
	}
	analysisOutput := getSingleVersionPackageOutput(analysis)
	output := struct {
		Image       string
		AnalyzeType string
		Analysis    []PackageOutput
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    analysisOutput,
	}
	return output
}

func (r SingleVersionPackageAnalyzeResult) OutputText(diffType string, format string) error {
	analysis, valid := r.Analysis.(map[string]PackageInfo)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type map[string]PackageInfo")
		return fmt.Errorf("Could not output %s analysis result", r.AnalyzeType)
	}
	analysisOutput := getSingleVersionPackageOutput(analysis)

	strAnalysis := stringifyPackages(analysisOutput)
	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    []StrPackageOutput
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    strAnalysis,
	}
	return TemplateOutputFromFormat(strResult, "SingleVersionPackageAnalyze", format)
}

type PackageOutput struct {
	Name    string
	Path    string `json:",omitempty"`
	Version string
	Size    int64
}

func getSingleVersionPackageOutput(packageMap map[string]PackageInfo) []PackageOutput {
	packages := []PackageOutput{}
	for name, info := range packageMap {
		packages = append(packages, PackageOutput{Name: name, Version: info.Version, Size: info.Size})
	}

	if SortSize {
		packageBy(packageSizeSort).Sort(packages)
	} else {
		packageBy(packageNameSort).Sort(packages)
	}
	return packages
}

func getMultiVersionPackageOutput(packageMap map[string]map[string]PackageInfo) []PackageOutput {
	packages := []PackageOutput{}
	for name, versionMap := range packageMap {
		for path, info := range versionMap {
			packages = append(packages, PackageOutput{Name: name, Path: path, Version: info.Version, Size: info.Size})
		}
	}

	if SortSize {
		packageBy(packageSizeSort).Sort(packages)
	} else {
		packageBy(packageNameSort).Sort(packages)
	}
	return packages
}

type FileAnalyzeResult AnalyzeResult

func (r FileAnalyzeResult) OutputStruct() interface{} {
	analysis, valid := r.Analysis.([]util.DirectoryEntry)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type []DirectoryEntry")
		return errors.New("Could not output FileAnalyzer analysis result")
	}

	if SortSize {
		directoryBy(directorySizeSort).Sort(analysis)
	} else {
		directoryBy(directoryNameSort).Sort(analysis)
	}
	r.Analysis = analysis
	return r
}

func (r FileAnalyzeResult) OutputText(analyzeType string, format string) error {
	analysis, valid := r.Analysis.([]util.DirectoryEntry)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type []DirectoryEntry")
		return errors.New("Could not output FileAnalyzer analysis result")
	}

	if SortSize {
		directoryBy(directorySizeSort).Sort(analysis)
	} else {
		directoryBy(directoryNameSort).Sort(analysis)
	}
	strAnalysis := stringifyDirectoryEntries(analysis)

	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    []StrDirectoryEntry
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    strAnalysis,
	}
	return TemplateOutputFromFormat(strResult, "FileAnalyze", format)
}
