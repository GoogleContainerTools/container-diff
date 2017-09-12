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
	"code.cloudfoundry.org/bytefmt"
	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
)

type StrPackageOutput struct {
	Name    string
	Path    string
	Version string
	Size    string
}

func stringifySize(size int64) string {
	strSize := "unknown"
	if size != -1 {
		strSize = bytefmt.ByteSize(uint64(size))
	}
	return strSize
}

func stringifyPackages(packages []PackageOutput) []StrPackageOutput {
	strPackages := []StrPackageOutput{}
	for _, pack := range packages {
		strSize := stringifySize(pack.Size)
		strPackages = append(strPackages, StrPackageOutput{pack.Name, pack.Path, pack.Version, strSize})
	}
	return strPackages
}

type StrMultiVersionInfo struct {
	Package string
	Info1   []StrPackageInfo
	Info2   []StrPackageInfo
}

type StrPackageInfo struct {
	Version string
	Size    string
}

func stringifyPackageInfo(info PackageInfo) StrPackageInfo {
	return StrPackageInfo{Version: info.Version, Size: stringifySize(info.Size)}
}

type StrInfo struct {
	Package string
	Info1   StrPackageInfo
	Info2   StrPackageInfo
}

func stringifyPackageDiff(infoDiff []Info) (strInfoDiff []StrInfo) {
	for _, diff := range infoDiff {
		strInfo1 := stringifyPackageInfo(diff.Info1)
		strInfo2 := stringifyPackageInfo(diff.Info2)

		strDiff := StrInfo{Package: diff.Package, Info1: strInfo1, Info2: strInfo2}
		strInfoDiff = append(strInfoDiff, strDiff)
	}
	return
}

func stringifyMultiVersionPackageDiff(infoDiff []MultiVersionInfo) (strInfoDiff []StrMultiVersionInfo) {
	for _, diff := range infoDiff {
		strInfos1 := []StrPackageInfo{}
		for _, info := range diff.Info1 {
			strInfos1 = append(strInfos1, stringifyPackageInfo(info))
		}

		strInfos2 := []StrPackageInfo{}
		for _, info := range diff.Info2 {
			strInfos2 = append(strInfos2, stringifyPackageInfo(info))
		}

		strDiff := StrMultiVersionInfo{Package: diff.Package, Info1: strInfos1, Info2: strInfos2}
		strInfoDiff = append(strInfoDiff, strDiff)
	}
	return
}

type StrDirectoryEntry struct {
	Name string
	Size string
}

func stringifyDirectoryEntries(entries []pkgutil.DirectoryEntry) (strEntries []StrDirectoryEntry) {
	for _, entry := range entries {
		strEntry := StrDirectoryEntry{Name: entry.Name, Size: stringifySize(entry.Size)}
		strEntries = append(strEntries, strEntry)
	}
	return
}

type StrEntryDiff struct {
	Name  string
	Size1 string
	Size2 string
}

func stringifyEntryDiffs(entries []EntryDiff) (strEntries []StrEntryDiff) {
	for _, entry := range entries {
		strEntry := StrEntryDiff{Name: entry.Name, Size1: stringifySize(entry.Size1), Size2: stringifySize(entry.Size2)}
		strEntries = append(strEntries, strEntry)
	}
	return
}
