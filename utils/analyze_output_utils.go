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
	Image string
	AnalyzeType string
	Analysis interface{}
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
		Image: r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis: analysisOutput,
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
		Image: r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis: analysisOutput,
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
	Name string
	Path string `json: ",omitempty"`
	Version string
	Size int64
}

type packageBy func(p1, p2 *PackageOutput) bool

func (by packageBy) Sort(packages []PackageOutput) {
	ps := &packageSorter{
		packages: packages,
		by: by,
	}
	sort.Sort(ps)
}

type packageSorter struct {
	packages []PackageOutput
	by func(p1, p2 *PackageOutput) bool
}

func (s *packageSorter) Len() int {
	return len(s.packages)
}

func (s *packageSorter) Less(i, j int) bool {
	return s.by(&s.packages[i], &s.packages[j])
}

func (s *packageSorter) Swap(i, j int) {
	s.packages[i], s.packages[j] = s.packages[i], s.packages[j]
}

var packageNameSort = func(p1, p2 *PackageOutput) bool {
	if p1.Name == p2.Name {
		return p1.Path < p2.Path
	}
	return p1.Name < p2.Name
}

var packageSizeSort = func(p1, p2 *PackageOutput) bool {
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
	Name string
	Path string
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
	return r
}

func (r FileAnalyzeResult) OutputText(analyzeType string) error {
	analysis := r.Analysis.([]DirectoryEntry)
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
