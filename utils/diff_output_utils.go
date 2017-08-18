package utils

import (
	"code.cloudfoundry.org/bytefmt"
)

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
	return TemplateOutput(r, "MultiVersionPackageDiff")
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
	return TemplateOutput(r, "SingleVersionPackageDiff")
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
		Image1 string
		Image2 string
		DiffType string
		Diff StrDiff
	}{
		Image1: r.Image1,
		Image2: r.Image2,
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
	Name string
	Size1 string
	Size2 string
}

func stringifyEntryDiffs(entries []EntryDiff) (strEntries []StrEntryDiff) {
	for _, entry := range entries {
		size1 := entry.Size1
		strSize1 := "unknown"
		if size1 != -1 {
			strSize1 = bytefmt.ByteSize(uint64(size1))
		}

		size2 := entry.Size2
		strSize2 := "unknown"
		if size2 != -1 {
			strSize2 = bytefmt.ByteSize(uint64(size2))
		}

		strEntry := StrEntryDiff{Name: entry.Name, Size1: strSize1, Size2: strSize2}
		strEntries = append(strEntries, strEntry)
	}
	return
}
