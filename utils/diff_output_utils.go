package utils

type DiffResult struct {
	Image1   string
	Image2   string
	DiffType string
	Diff     interface{}
}

type MultiVersionPackageDiffResult DiffResult

func (r MultiVersionPackageDiffResult) OutputStruct() interface{} {
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

type SingleVersionPackageDiffResult DiffResult

func (r SingleVersionPackageDiffResult) OutputStruct() interface{} {
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

type HistDiffResult DiffResult

func (r HistDiffResult) OutputStruct() interface{} {
	return r
}

func (r HistDiffResult) OutputText(diffType string) error {
	return TemplateOutput(r, "HistDiff")
}

type DirDiffResult DiffResult

func (r DirDiffResult) OutputStruct() interface{} {
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
