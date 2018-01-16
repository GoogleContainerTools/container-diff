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

func (r MultiVersionPackageDiffResult) OutputText(diffType string, format string) error {
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
	return TemplateOutputFromFormat(strResult, "MultiVersionPackageDiff", format)
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

func (r SingleVersionPackageDiffResult) OutputText(diffType string, format string) error {
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
	return TemplateOutputFromFormat(strResult, "SingleVersionPackageDiff", format)
}

func getSingleVersionInfoDiffOutput(infoDiff []Info) []Info {
	if SortSize {
		singleInfoBy(singleInfoSizeSort).Sort(infoDiff)
	} else {
		singleInfoBy(singleInfoNameSort).Sort(infoDiff)
	}
	return infoDiff
}

type HistDiffResult DiffResult

func (r HistDiffResult) OutputStruct() interface{} {
	return r
}

func (r HistDiffResult) OutputText(diffType string, format string) error {
	return TemplateOutputFromFormat(r, "HistDiff", format)
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

func (r DirDiffResult) OutputText(diffType string, format string) error {
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
	return TemplateOutputFromFormat(strResult, "DirDiff", format)
}
