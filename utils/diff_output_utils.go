package utils

import (
	"sort"
)

type DiffResult struct {
	Image1   string
	Image2   string
	DiffType string
	Diff     interface{}
}

type MultiVersionPackageDiffResult DiffResult

func (r MultiVersionPackageDiffResult) GetStruct() interface{} {
	diff := r.Diff.(MultiVersionPackageDiff)
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

func (r MultiVersionPackageDiffResult) OutputText(diffType string) error {
	diff := r.Diff.(MultiVersionPackageDiff)

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
	return TemplateOutput(strResult, "MultiVersionPackageDiff")
}

func getMultiVersionInfoDiffOutput(infoDiff []MultiVersionInfo) []MultiVersionInfo {
	if SortSize {
		multiInfoBy(multiInfoSizeSort).Sort(infoDiff)
	} else {
		multiInfoBy(multiInfoNameSort).Sort(infoDiff)
	}
	return infoDiff
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
	s.packageDiffs[i], s.packageDiffs[j] = s.packageDiffs[i], s.packageDiffs[j]
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

type SingleVersionPackageDiffResult DiffResult

func (r SingleVersionPackageDiffResult) GetStruct() interface{} {
	diff := r.Diff.(PackageDiff)
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

func (r SingleVersionPackageDiffResult) OutputText(diffType string) error {
	diff := r.Diff.(PackageDiff)

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
	return TemplateOutput(strResult, "SingleVersionPackageDiff")
}

func getSingleVersionInfoDiffOutput(infoDiff []Info) []Info {
	if SortSize {
		singleInfoBy(singleInfoSizeSort).Sort(infoDiff)
	} else {
		singleInfoBy(singleInfoNameSort).Sort(infoDiff)
	}
	return infoDiff
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
	s.packageDiffs[i], s.packageDiffs[j] = s.packageDiffs[i], s.packageDiffs[j]
}

var singleInfoNameSort = func(a, b *Info) bool {
	return a.Package < b.Package
}

// Sorts MultiVersionInfos by package instance with the largest size in the first image, in descending order
var singleInfoSizeSort = func(a, b *Info) bool {
	aInfo1 := a.Info1
	bInfo1 := b.Info1

	// Compares the sizes of the packages in the first image
	return aInfo1.Size > bInfo1.Size
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

type HistDiffResult DiffResult

func (r HistDiffResult) GetStruct() interface{} {
	return r
}

func (r HistDiffResult) OutputText(diffType string) error {
	return TemplateOutput(r, "HistDiff")
}

type DirDiffResult DiffResult

func (r DirDiffResult) GetStruct() interface{} {
	r.Diff = sortDirDiff(r.Diff.(DirDiff))
	return r
}

func (r DirDiffResult) OutputText(diffType string) error {
	diff := sortDirDiff(r.Diff.(DirDiff))

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
	return TemplateOutput(strResult, "DirDiff")
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
	s.entryDiffs[i], s.entryDiffs[j] = s.entryDiffs[i], s.entryDiffs[j]
}

var entryDiffNameSort = func(a, b *EntryDiff) bool {
	return a.Name < b.Name
}

// Sorts by size of the files in the first image, in descending order
var entryDiffSizeSort = func(a, b *EntryDiff) bool {
	return a.Size1 > b.Size1
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
