package utils

import (
	"sort"

	"code.cloudfoundry.org/bytefmt"
)

var SortSize bool

type Result interface {
	GetStruct() interface{}
	OutputText(resultType string) error
}

type AnalyzeResult struct {
	Image       string
	AnalyzeType string
	Analysis    interface{}
}

type ListAnalyzeResult AnalyzeResult

func (r ListAnalyzeResult) GetStruct() interface{} {
	return r
}

func (r ListAnalyzeResult) OutputText(resultType string) error {
	r.Analysis = r.Analysis.([]string)
	return TemplateOutput(r, "ListAnalyze")
}

type MultiVersionPackageAnalyzeResult AnalyzeResult

func (r MultiVersionPackageAnalyzeResult) GetStruct() interface{} {
	analysis := r.Analysis.(map[string]map[string]PackageInfo)
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

func (r MultiVersionPackageAnalyzeResult) OutputText(resultType string) error {
	analysis := r.Analysis.(map[string]map[string]PackageInfo)
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
	return TemplateOutput(strResult, "MultiVersionPackageAnalyze")
}

type SingleVersionPackageAnalyzeResult AnalyzeResult

func (r SingleVersionPackageAnalyzeResult) GetStruct() interface{} {
	analysis := r.Analysis.(map[string]PackageInfo)
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

func (r SingleVersionPackageAnalyzeResult) OutputText(diffType string) error {
	analysis := r.Analysis.(map[string]PackageInfo)
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
	return TemplateOutput(strResult, "SingleVersionPackageAnalyze")
}

type PackageOutput struct {
	Name    string
	Path    string `json:",omitempty"`
	Version string
	Size    int64
}

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

type FileAnalyzeResult AnalyzeResult

func (r FileAnalyzeResult) GetStruct() interface{} {
	analysis := r.Analysis.([]DirectoryEntry)
	if SortSize {
		directoryBy(directorySizeSort).Sort(analysis)
	} else {
		directoryBy(directoryNameSort).Sort(analysis)
	}
	r.Analysis = analysis
	return r
}

func (r FileAnalyzeResult) OutputText(analyzeType string) error {
	analysis := r.Analysis.([]DirectoryEntry)
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
	return TemplateOutput(strResult, "FileAnalyze")
}

type directoryBy func(e1, e2 *DirectoryEntry) bool

func (by directoryBy) Sort(entries []DirectoryEntry) {
	ds := &directorySorter{
		entries: entries,
		by:      by,
	}
	sort.Sort(ds)
}

type directorySorter struct {
	entries []DirectoryEntry
	by      func(p1, p2 *DirectoryEntry) bool
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

// If packages have the same name, means they exist where multiple version of the same package are allowed,
// so sort by version.  If they have the same version, then sort by size.
var directoryNameSort = func(e1, e2 *DirectoryEntry) bool {
	return e1.Name < e2.Name
}

// If packages have the same size, sort by name.  If they are two versions of the same package, sort by version.
var directorySizeSort = func(e1, e2 *DirectoryEntry) bool {
	if e1.Size == e2.Size {
		return e1.Name < e2.Name
	}
	return e1.Size > e2.Size
}

type StrDirectoryEntry struct {
	Name string
	Size string
}

func stringifyDirectoryEntries(entries []DirectoryEntry) (strEntries []StrDirectoryEntry) {
	for _, entry := range entries {
		strEntry := StrDirectoryEntry{Name: entry.Name, Size: stringifySize(entry.Size)}
		strEntries = append(strEntries, strEntry)
	}
	return
}
