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

package util

import (
	"errors"
	"fmt"
	"io"

	"github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/sirupsen/logrus"
)

type Result interface {
	OutputStruct() interface{}
	OutputText(writer io.Writer, resultType string, format string) error
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

func (r ListAnalyzeResult) OutputText(writer io.Writer, resultType string, format string) error {
	analysis, valid := r.Analysis.([]string)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type []string")
		return fmt.Errorf("Could not output %s analysis result", r.AnalyzeType)
	}
	r.Analysis = analysis
	return TemplateOutputFromFormat(writer, r, "ListAnalyze", format)

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

func (r MultiVersionPackageAnalyzeResult) OutputText(writer io.Writer, resultType string, format string) error {
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
	return TemplateOutputFromFormat(writer, strResult, "MultiVersionPackageAnalyze", format)
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

func (r SingleVersionPackageAnalyzeResult) OutputText(writer io.Writer, diffType string, format string) error {
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
	return TemplateOutputFromFormat(writer, strResult, "SingleVersionPackageAnalyze", format)
}

type SingleVersionPackageLayerAnalyzeResult AnalyzeResult

func (r SingleVersionPackageLayerAnalyzeResult) OutputStruct() interface{} {
	analysis, valid := r.Analysis.(PackageLayerDiff)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type PackageLayerDiff")
		return fmt.Errorf("Could not output %s analysis result", r.AnalyzeType)
	}

	type PkgDiff struct {
		Packages1 []PackageOutput
		Packages2 []PackageOutput
		InfoDiff  []Info
	}

	var analysisOutput []PkgDiff
	for _, d := range analysis.PackageDiffs {
		diffOutput := PkgDiff{
			Packages1: getSingleVersionPackageOutput(d.Packages1),
			Packages2: getSingleVersionPackageOutput(d.Packages2),
			InfoDiff:  getSingleVersionInfoDiffOutput(d.InfoDiff),
		}
		analysisOutput = append(analysisOutput, diffOutput)
	}

	output := struct {
		Image       string
		AnalyzeType string
		Analysis    []PkgDiff
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    analysisOutput,
	}
	return output
}

func (r SingleVersionPackageLayerAnalyzeResult) OutputText(writer io.Writer, diffType string, format string) error {
	analysis, valid := r.Analysis.(PackageLayerDiff)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type PackageLayerDiff")
		return fmt.Errorf("Could not output %s analysis result", r.AnalyzeType)
	}

	type StrDiff struct {
		Packages1 []StrPackageOutput
		Packages2 []StrPackageOutput
		InfoDiff  []StrInfo
	}

	var analysisOutput []StrDiff
	for _, d := range analysis.PackageDiffs {
		diffOutput := StrDiff{
			Packages1: stringifyPackages(getSingleVersionPackageOutput(d.Packages1)),
			Packages2: stringifyPackages(getSingleVersionPackageOutput(d.Packages2)),
			InfoDiff:  stringifyPackageDiff(getSingleVersionInfoDiffOutput(d.InfoDiff)),
		}
		analysisOutput = append(analysisOutput, diffOutput)
	}

	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    []StrDiff
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    analysisOutput,
	}
	return TemplateOutputFromFormat(writer, strResult, "SingleVersionPackageLayerAnalyze", format)
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

func (r FileAnalyzeResult) OutputText(writer io.Writer, analyzeType string, format string) error {
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
	return TemplateOutputFromFormat(writer, strResult, "FileAnalyze", format)
}

type FileLayerAnalyzeResult AnalyzeResult

func (r FileLayerAnalyzeResult) OutputStruct() interface{} {
	analysis, valid := r.Analysis.([][]util.DirectoryEntry)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type []DirectoryEntry")
		return errors.New("Could not output FileAnalyzer analysis result")
	}

	for _, a := range analysis {
		if SortSize {
			directoryBy(directorySizeSort).Sort(a)
		} else {
			directoryBy(directoryNameSort).Sort(a)
		}
	}

	r.Analysis = analysis
	return r
}

func (r FileLayerAnalyzeResult) OutputText(writer io.Writer, analyzeType string, format string) error {
	analysis, valid := r.Analysis.([][]util.DirectoryEntry)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type []DirectoryEntry")
		return errors.New("Could not output FileAnalyzer analysis result")
	}

	var strDirectoryEntries [][]StrDirectoryEntry

	for _, a := range analysis {
		if SortSize {
			directoryBy(directorySizeSort).Sort(a)
		} else {
			directoryBy(directoryNameSort).Sort(a)
		}
		strAnalysis := stringifyDirectoryEntries(a)
		strDirectoryEntries = append(strDirectoryEntries, strAnalysis)
	}

	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    [][]StrDirectoryEntry
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    strDirectoryEntries,
	}
	return TemplateOutputFromFormat(writer, strResult, "FileLayerAnalyze", format)
}

type SizeAnalyzeResult AnalyzeResult

func (r SizeAnalyzeResult) OutputStruct() interface{} {
	analysis, valid := r.Analysis.([]SizeEntry)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type []SizeEntry")
		return errors.New("Could not output SizeAnalyzer analysis result")
	}
	r.Analysis = analysis
	return r
}

func (r SizeAnalyzeResult) OutputText(writer io.Writer, analyzeType string, format string) error {
	analysis, valid := r.Analysis.([]SizeEntry)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type []SizeEntry")
		return errors.New("Could not output SizeAnalyzer analysis result")
	}

	strAnalysis := stringifySizeEntries(analysis)

	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    []StrSizeEntry
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    strAnalysis,
	}
	return TemplateOutputFromFormat(writer, strResult, "SizeAnalyze", format)
}

type SizeLayerAnalyzeResult AnalyzeResult

func (r SizeLayerAnalyzeResult) OutputStruct() interface{} {
	analysis, valid := r.Analysis.([]SizeEntry)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type []SizeEntry")
		return errors.New("Could not output SizeLayerAnalyzer analysis result")
	}
	r.Analysis = analysis
	return r
}

func (r SizeLayerAnalyzeResult) OutputText(writer io.Writer, analyzeType string, format string) error {
	analysis, valid := r.Analysis.([]SizeEntry)
	if !valid {
		logrus.Error("Unexpected structure of Analysis.  Should be of type []SizeEntry")
		return errors.New("Could not output SizeLayerAnalyzer analysis result")
	}

	strAnalysis := stringifySizeEntries(analysis)

	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    []StrSizeEntry
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    strAnalysis,
	}
	return TemplateOutputFromFormat(writer, strResult, "SizeLayerAnalyze", format)
}
