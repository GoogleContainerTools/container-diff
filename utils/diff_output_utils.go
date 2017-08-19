package utils

type DiffResult interface {
	GetStruct() DiffResult
	OutputText(diffType string) error
}

type MultiVersionPackageDiffResult struct {
	Image1   string
	Image2   string
	DiffType string
	Diff     MultiVersionPackageDiff
}

func (r MultiVersionPackageDiffResult) GetStruct() DiffResult {
	return r
}

func (r MultiVersionPackageDiffResult) OutputText(diffType string) error {
	diff := r.Diff

	strPackages1 := stringifyMultiVersionPackages(diff.Packages1)
	strPackages2 := stringifyMultiVersionPackages(diff.Packages2)
	strInfoDiff := stringifyMultiVersionPackageDiff(diff.InfoDiff)

	type StrDiff struct {
		Packages1 map[string]map[string]StrPackageInfo
		Packages2 map[string]map[string]StrPackageInfo
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

type StrMultiVersionInfo struct {
	Package string
	Info1   []StrPackageInfo
	Info2   []StrPackageInfo
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

type SingleVersionPackageDiffResult struct {
	Image1   string
	Image2   string
	DiffType string
	Diff     PackageDiff
}

func (r SingleVersionPackageDiffResult) GetStruct() DiffResult {
	return r
}

func (r SingleVersionPackageDiffResult) OutputText(diffType string) error {
	diff := r.Diff

	strPackages1 := stringifyPackages(diff.Packages1)
	strPackages2 := stringifyPackages(diff.Packages2)
	strInfoDiff := stringifyPackageDiff(diff.InfoDiff)

	type StrDiff struct {
		Packages1 map[string]StrPackageInfo
		Packages2 map[string]StrPackageInfo
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

type HistDiffResult struct {
	Image1   string
	Image2   string
	DiffType string
	Diff     HistDiff
}

func (r HistDiffResult) GetStruct() DiffResult {
	return r
}

func (r HistDiffResult) OutputText(diffType string) error {
	return TemplateOutput(r, "HistDiff")
}

type DirDiffResult struct {
	Image1   string
	Image2   string
	DiffType string
	Diff     DirDiff
}

func (r DirDiffResult) GetStruct() DiffResult {
	return r
}

func (r DirDiffResult) OutputText(diffType string) error {
	diff := r.Diff

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
