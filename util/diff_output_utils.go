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

	"github.com/sirupsen/logrus"
)

type DiffResult struct {
	Image1   string
	Image2   string
	DiffType string
	Diff     interface{}
}

type MultiVersionPackageDiffResult DiffResult

func (r MultiVersionPackageDiffResult) OutputStruct() interface{} {
	diff, valid := r.Diff.(MultiVersionPackageDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should follow the MultiVersionPackageDiff struct")
		return fmt.Errorf("Could not output %s diff result", r.DiffType)
	}

	diffOutput := struct {
		Packages1 []PackageOutput
		Packages2 []PackageOutput
		InfoDiff  []MultiVersionInfo
	}{
		Packages1: getMultiVersionPackageOutput(diff.Packages1),
		Packages2: getMultiVersionPackageOutput(diff.Packages2),
		InfoDiff:  getMultiVersionInfoDiffOutput(diff.InfoDiff),
	}
	r.Diff = diffOutput
	return r
}

func (r MultiVersionPackageDiffResult) OutputText(writer io.Writer, diffType string, format string) error {
	diff, valid := r.Diff.(MultiVersionPackageDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should follow the MultiVersionPackageDiff struct")
		return fmt.Errorf("Could not output %s diff result", r.DiffType)
	}

	strPackages1 := stringifyPackages(getMultiVersionPackageOutput(diff.Packages1))
	strPackages2 := stringifyPackages(getMultiVersionPackageOutput(diff.Packages2))
	strInfoDiff := stringifyMultiVersionPackageDiff(getMultiVersionInfoDiffOutput(diff.InfoDiff))

	type StrDiff struct {
		Packages1 []StrPackageOutput
		Packages2 []StrPackageOutput
		InfoDiff  []StrMultiVersionInfo
	}

	strResult := struct {
		Image1   string
		Image2   string
		DiffType string
		Diff     StrDiff
	}{
		Image1:   r.Image1,
		Image2:   r.Image2,
		DiffType: r.DiffType,
		Diff: StrDiff{
			Packages1: strPackages1,
			Packages2: strPackages2,
			InfoDiff:  strInfoDiff,
		},
	}
	return TemplateOutputFromFormat(writer, strResult, "MultiVersionPackageDiff", format)
}

func getMultiVersionInfoDiffOutput(infoDiff []MultiVersionInfo) []MultiVersionInfo {
	if SortSize {
		multiInfoBy(multiInfoSizeSort).Sort(infoDiff)
	} else {
		multiInfoBy(multiInfoNameSort).Sort(infoDiff)
	}
	return infoDiff
}

type SingleVersionPackageDiffResult DiffResult

func (r SingleVersionPackageDiffResult) OutputStruct() interface{} {
	diff, valid := r.Diff.(PackageDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should follow the PackageDiff struct")
		return fmt.Errorf("Could not output %s diff result", r.DiffType)
	}

	diffOutput := struct {
		Packages1 []PackageOutput
		Packages2 []PackageOutput
		InfoDiff  []Info
	}{
		Packages1: getSingleVersionPackageOutput(diff.Packages1),
		Packages2: getSingleVersionPackageOutput(diff.Packages2),
		InfoDiff:  getSingleVersionInfoDiffOutput(diff.InfoDiff),
	}
	r.Diff = diffOutput
	return r
}

func (r SingleVersionPackageDiffResult) OutputText(writer io.Writer, diffType string, format string) error {
	diff, valid := r.Diff.(PackageDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should follow the PackageDiff struct")
		return fmt.Errorf("Could not output %s diff result", r.DiffType)
	}

	strPackages1 := stringifyPackages(getSingleVersionPackageOutput(diff.Packages1))
	strPackages2 := stringifyPackages(getSingleVersionPackageOutput(diff.Packages2))
	strInfoDiff := stringifyPackageDiff(getSingleVersionInfoDiffOutput(diff.InfoDiff))

	type StrDiff struct {
		Packages1 []StrPackageOutput
		Packages2 []StrPackageOutput
		InfoDiff  []StrInfo
	}

	strResult := struct {
		Image1   string
		Image2   string
		DiffType string
		Diff     StrDiff
	}{
		Image1:   r.Image1,
		Image2:   r.Image2,
		DiffType: r.DiffType,
		Diff: StrDiff{
			Packages1: strPackages1,
			Packages2: strPackages2,
			InfoDiff:  strInfoDiff,
		},
	}
	return TemplateOutputFromFormat(writer, strResult, "SingleVersionPackageDiff", format)
}

func getSingleVersionInfoDiffOutput(infoDiff []Info) []Info {
	if SortSize {
		singleInfoBy(singleInfoSizeSort).Sort(infoDiff)
	} else {
		singleInfoBy(singleInfoNameSort).Sort(infoDiff)
	}
	return infoDiff
}

type SingleVersionPackageLayerDiffResult DiffResult

func (r SingleVersionPackageLayerDiffResult) OutputStruct() interface{} {
	diff, valid := r.Diff.(PackageLayerDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should follow the PackageLayerDiff struct")
		return fmt.Errorf("Could not output %s diff result", r.DiffType)
	}

	type PkgDiff struct {
		Packages1 []PackageOutput
		Packages2 []PackageOutput
		InfoDiff  []Info
	}

	var diffOutputs []PkgDiff
	for _, d := range diff.PackageDiffs {
		diffOutput := PkgDiff{
			Packages1: getSingleVersionPackageOutput(d.Packages1),
			Packages2: getSingleVersionPackageOutput(d.Packages2),
			InfoDiff:  getSingleVersionInfoDiffOutput(d.InfoDiff),
		}
		diffOutputs = append(diffOutputs, diffOutput)
	}

	r.Diff = diffOutputs
	return r
}

func (r SingleVersionPackageLayerDiffResult) OutputText(writer io.Writer, diffType string, format string) error {
	diff, valid := r.Diff.(PackageLayerDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should follow the PackageLayerDiff struct")
		return fmt.Errorf("Could not output %s diff result", r.DiffType)
	}

	type StrDiff struct {
		Packages1 []StrPackageOutput
		Packages2 []StrPackageOutput
		InfoDiff  []StrInfo
	}

	var diffOutputs []StrDiff
	for _, d := range diff.PackageDiffs {
		diffOutput := StrDiff{
			Packages1: stringifyPackages(getSingleVersionPackageOutput(d.Packages1)),
			Packages2: stringifyPackages(getSingleVersionPackageOutput(d.Packages2)),
			InfoDiff:  stringifyPackageDiff(getSingleVersionInfoDiffOutput(d.InfoDiff)),
		}
		diffOutputs = append(diffOutputs, diffOutput)
	}

	strResult := struct {
		Image1   string
		Image2   string
		DiffType string
		Diff     []StrDiff
	}{
		Image1:   r.Image1,
		Image2:   r.Image2,
		DiffType: r.DiffType,
		Diff:     diffOutputs,
	}
	return TemplateOutputFromFormat(writer, strResult, "SingleVersionPackageLayerDiff", format)
}

type HistDiffResult DiffResult

func (r HistDiffResult) OutputStruct() interface{} {
	return r
}

func (r HistDiffResult) OutputText(writer io.Writer, diffType string, format string) error {
	return TemplateOutputFromFormat(writer, r, "HistDiff", format)
}

type MetadataDiffResult DiffResult

func (r MetadataDiffResult) OutputStruct() interface{} {
	return r
}

func (r MetadataDiffResult) OutputText(writer io.Writer, diffType string, format string) error {
	return TemplateOutputFromFormat(writer, r, "MetadataDiff", format)
}

type DirDiffResult DiffResult

func (r DirDiffResult) OutputStruct() interface{} {
	diff, valid := r.Diff.(DirDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should follow the DirDiff struct")
		return errors.New("Could not output FileAnalyzer diff result")
	}

	r.Diff = sortDirDiff(diff)
	return r
}

func (r DirDiffResult) OutputText(writer io.Writer, diffType string, format string) error {
	diff, valid := r.Diff.(DirDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should follow the DirDiff struct")
		return errors.New("Could not output FileAnalyzer diff result")
	}
	diff = sortDirDiff(diff)

	strAdds := stringifyDirectoryEntries(diff.Adds)
	strDels := stringifyDirectoryEntries(diff.Dels)
	strMods := stringifyEntryDiffs(diff.Mods)

	type StrDiff struct {
		Adds []StrDirectoryEntry
		Dels []StrDirectoryEntry
		Mods []StrEntryDiff
	}

	strResult := struct {
		Image1   string
		Image2   string
		DiffType string
		Diff     StrDiff
	}{
		Image1:   r.Image1,
		Image2:   r.Image2,
		DiffType: r.DiffType,
		Diff: StrDiff{
			Adds: strAdds,
			Dels: strDels,
			Mods: strMods,
		},
	}
	return TemplateOutputFromFormat(writer, strResult, "DirDiff", format)
}

type SizeDiffResult DiffResult

func (r SizeDiffResult) OutputStruct() interface{} {
	diff, valid := r.Diff.([]SizeDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should be of type []SizeDiff")
		return errors.New("Could not output SizeAnalyzer diff result")
	}

	r.Diff = diff
	return r
}

func (r SizeDiffResult) OutputText(writer io.Writer, diffType string, format string) error {
	diff, valid := r.Diff.([]SizeDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should be of type []SizeDiff")
		return errors.New("Could not output SizeAnalyzer diff result")
	}

	strDiff := stringifySizeDiffs(diff)

	strResult := struct {
		Image1   string
		Image2   string
		DiffType string
		Diff     []StrSizeDiff
	}{
		Image1:   r.Image1,
		Image2:   r.Image2,
		DiffType: r.DiffType,
		Diff:     strDiff,
	}
	return TemplateOutputFromFormat(writer, strResult, "SizeDiff", format)
}

type SizeLayerDiffResult DiffResult

func (r SizeLayerDiffResult) OutputStruct() interface{} {
	diff, valid := r.Diff.([]SizeDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should be of type []SizeDiff")
		return errors.New("Could not output SizeLayerAnalyzer diff result")
	}

	r.Diff = diff
	return r
}

func (r SizeLayerDiffResult) OutputText(writer io.Writer, diffType string, format string) error {
	diff, valid := r.Diff.([]SizeDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should be of type []SizeDiff")
		return errors.New("Could not output SizeLayerAnalyzer diff result")
	}

	strDiff := stringifySizeDiffs(diff)

	strResult := struct {
		Image1   string
		Image2   string
		DiffType string
		Diff     []StrSizeDiff
	}{
		Image1:   r.Image1,
		Image2:   r.Image2,
		DiffType: r.DiffType,
		Diff:     strDiff,
	}
	return TemplateOutputFromFormat(writer, strResult, "SizeLayerDiff", format)
}

type MultipleDirDiffResult DiffResult

func (r MultipleDirDiffResult) OutputStruct() interface{} {
	diff, valid := r.Diff.(MultipleDirDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should follow the MultipleDirDiff struct")
		return errors.New("Could not output FileLayerAnalyzer diff result")
	}
	for i, d := range diff.DirDiffs {
		diff.DirDiffs[i] = sortDirDiff(d)
	}
	r.Diff = diff
	return r
}

func (r MultipleDirDiffResult) OutputText(writer io.Writer, diffType string, format string) error {
	diff, valid := r.Diff.(MultipleDirDiff)
	if !valid {
		logrus.Error("Unexpected structure of Diff.  Should follow the MultipleDirDiff struct")
		return errors.New("Could not output FileLayerAnalyzer diff result")
	}
	for i, d := range diff.DirDiffs {
		diff.DirDiffs[i] = sortDirDiff(d)
	}

	type StrDiff struct {
		Adds []StrDirectoryEntry
		Dels []StrDirectoryEntry
		Mods []StrEntryDiff
	}

	var strDiffs []StrDiff
	for _, d := range diff.DirDiffs {
		strAdds := stringifyDirectoryEntries(d.Adds)
		strDels := stringifyDirectoryEntries(d.Dels)
		strMods := stringifyEntryDiffs(d.Mods)

		strDiffs = append(strDiffs, StrDiff{
			Adds: strAdds,
			Dels: strDels,
			Mods: strMods,
		})

	}

	type ImageDiff struct {
		StrDiffs []StrDiff
	}
	strResult := struct {
		Image1   string
		Image2   string
		DiffType string
		Diff     []StrDiff
	}{
		Image1:   r.Image1,
		Image2:   r.Image2,
		DiffType: r.DiffType,
		Diff:     strDiffs,
	}
	return TemplateOutputFromFormat(writer, strResult, "MultipleDirDiff", format)
}
