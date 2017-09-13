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
	"sort"

	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
)

var SortSize bool

type packageBy func(p1, p2 *PackageOutput) bool

func (by packageBy) Sort(packages []PackageOutput) {
	ps := &packageSorter{
		packages: packages,
		by:       by,
	}
	sort.Sort(ps)
}

type packageSorter struct {
	packages []PackageOutput
	by       func(p1, p2 *PackageOutput) bool
}

func (s *packageSorter) Len() int {
	return len(s.packages)
}

func (s *packageSorter) Less(i, j int) bool {
	return s.by(&s.packages[i], &s.packages[j])
}

func (s *packageSorter) Swap(i, j int) {
	s.packages[i], s.packages[j] = s.packages[j], s.packages[i]
}

// If packages have the same name, means they exist where multiple version of the same package are allowed,
// so sort by version.  If they have the same version, then sort by size.
var packageNameSort = func(p1, p2 *PackageOutput) bool {
	if p1.Name == p2.Name {
		if p1.Version == p2.Version {
			return p1.Size > p2.Size
		}
		return p1.Version < p2.Version
	}
	return p1.Name < p2.Name
}

// If packages have the same size, sort by name.  If they are two versions of the same package, sort by version.
var packageSizeSort = func(p1, p2 *PackageOutput) bool {
	if p1.Size == p2.Size {
		if p1.Name == p2.Name {
			return p1.Version < p2.Version
		}
		return p1.Name < p2.Name
	}
	return p1.Size > p2.Size
}

type singleInfoBy func(a, b *Info) bool

func (by singleInfoBy) Sort(packageDiffs []Info) {
	ss := &singleVersionInfoSorter{
		packageDiffs: packageDiffs,
		by:           by,
	}
	sort.Sort(ss)
}

type singleVersionInfoSorter struct {
	packageDiffs []Info
	by           func(a, b *Info) bool
}

func (s *singleVersionInfoSorter) Len() int {
	return len(s.packageDiffs)
}

func (s *singleVersionInfoSorter) Less(i, j int) bool {
	return s.by(&s.packageDiffs[i], &s.packageDiffs[j])
}

func (s *singleVersionInfoSorter) Swap(i, j int) {
	s.packageDiffs[i], s.packageDiffs[j] = s.packageDiffs[j], s.packageDiffs[i]
}

var singleInfoNameSort = func(a, b *Info) bool {
	return a.Package < b.Package
}

// Sorts MultiVersionInfos by package instance with the largest size in the first image, in descending order
var singleInfoSizeSort = func(a, b *Info) bool {
	return a.Info1.Size > b.Info1.Size
}

type multiInfoBy func(a, b *MultiVersionInfo) bool

func (by multiInfoBy) Sort(packageDiffs []MultiVersionInfo) {
	ms := &multiVersionInfoSorter{
		packageDiffs: packageDiffs,
		by:           by,
	}
	sort.Sort(ms)
}

type multiVersionInfoSorter struct {
	packageDiffs []MultiVersionInfo
	by           func(a, b *MultiVersionInfo) bool
}

func (s *multiVersionInfoSorter) Len() int {
	return len(s.packageDiffs)
}

func (s *multiVersionInfoSorter) Less(i, j int) bool {
	return s.by(&s.packageDiffs[i], &s.packageDiffs[j])
}

func (s *multiVersionInfoSorter) Swap(i, j int) {
	s.packageDiffs[i], s.packageDiffs[j] = s.packageDiffs[j], s.packageDiffs[i]
}

var multiInfoNameSort = func(a, b *MultiVersionInfo) bool {
	return a.Package < b.Package
}

// Sorts MultiVersionInfos by package instance with the largest size in the first image, in descending order
var multiInfoSizeSort = func(a, b *MultiVersionInfo) bool {
	aInfo1 := a.Info1
	bInfo1 := b.Info1

	// For each package, sorts the infos of the first image's instances of that package in descending order
	sort.Sort(packageInfoBySize(aInfo1))
	sort.Sort(packageInfoBySize(bInfo1))
	// Compares the largest size instances of each package in the first image
	return aInfo1[0].Size > bInfo1[0].Size
}

type packageInfoBySize []PackageInfo

func (infos packageInfoBySize) Len() int {
	return len(infos)
}

func (infos packageInfoBySize) Swap(i, j int) {
	infos[i], infos[j] = infos[j], infos[i]
}

func (infos packageInfoBySize) Less(i, j int) bool {
	if infos[i].Size == infos[j].Size {
		return infos[i].Version < infos[j].Version
	}
	return infos[i].Size > infos[j].Size
}

type packageInfoByVersion []PackageInfo

func (infos packageInfoByVersion) Len() int {
	return len(infos)
}

func (infos packageInfoByVersion) Swap(i, j int) {
	infos[i], infos[j] = infos[j], infos[i]
}

func (infos packageInfoByVersion) Less(i, j int) bool {
	if infos[i].Version == infos[j].Version {
		return infos[i].Size > infos[j].Size
	}
	return infos[i].Version < infos[j].Version
}

type directoryBy func(e1, e2 *pkgutil.DirectoryEntry) bool

func (by directoryBy) Sort(entries []pkgutil.DirectoryEntry) {
	ds := &directorySorter{
		entries: entries,
		by:      by,
	}
	sort.Sort(ds)
}

type directorySorter struct {
	entries []pkgutil.DirectoryEntry
	by      func(p1, p2 *pkgutil.DirectoryEntry) bool
}

func (s *directorySorter) Len() int {
	return len(s.entries)
}

func (s *directorySorter) Less(i, j int) bool {
	return s.by(&s.entries[i], &s.entries[j])
}

func (s *directorySorter) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
}

var directoryNameSort = func(e1, e2 *pkgutil.DirectoryEntry) bool {
	return e1.Name < e2.Name
}

// If directory entries have the same size, sort by name.
var directorySizeSort = func(e1, e2 *pkgutil.DirectoryEntry) bool {
	if e1.Size == e2.Size {
		return e1.Name < e2.Name
	}
	return e1.Size > e2.Size
}

func sortDirDiff(diff DirDiff) DirDiff {
	adds, dels, mods := diff.Adds, diff.Dels, diff.Mods
	if SortSize {
		directoryBy(directorySizeSort).Sort(adds)
		directoryBy(directorySizeSort).Sort(dels)
		entryDiffBy(entryDiffSizeSort).Sort(mods)
	} else {
		directoryBy(directoryNameSort).Sort(adds)
		directoryBy(directoryNameSort).Sort(dels)
		entryDiffBy(entryDiffSizeSort).Sort(mods)
	}
	return DirDiff{adds, dels, mods}
}

type entryDiffBy func(a, b *EntryDiff) bool

func (by entryDiffBy) Sort(entryDiffs []EntryDiff) {
	ds := &entryDiffSorter{
		entryDiffs: entryDiffs,
		by:         by,
	}
	sort.Sort(ds)
}

type entryDiffSorter struct {
	entryDiffs []EntryDiff
	by         func(a, b *EntryDiff) bool
}

func (s *entryDiffSorter) Len() int {
	return len(s.entryDiffs)
}

func (s *entryDiffSorter) Less(i, j int) bool {
	return s.by(&s.entryDiffs[i], &s.entryDiffs[j])
}

func (s *entryDiffSorter) Swap(i, j int) {
	s.entryDiffs[i], s.entryDiffs[j] = s.entryDiffs[j], s.entryDiffs[i]
}

var entryDiffNameSort = func(a, b *EntryDiff) bool {
	return a.Name < b.Name
}

// Sorts by size of the files in the first image, in descending order
var entryDiffSizeSort = func(a, b *EntryDiff) bool {
	return a.Size1 > b.Size1
}
