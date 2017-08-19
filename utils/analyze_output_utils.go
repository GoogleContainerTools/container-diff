package utils

import (
	"code.cloudfoundry.org/bytefmt"
)

type AnalyzeResult interface {
	GetStruct() AnalyzeResult
	OutputText(analyzeType string) error
}

type ListAnalyzeResult struct {
	Image       string
	AnalyzeType string
	Analysis    []string
}

func (r ListAnalyzeResult) GetStruct() AnalyzeResult {
	return r
}

func (r ListAnalyzeResult) OutputText(analyzeType string) error {
	return TemplateOutput(r, "ListAnalyze")
}

type MultiVersionPackageAnalyzeResult struct {
	Image       string
	AnalyzeType string
	Analysis    map[string]map[string]PackageInfo
}

func (r MultiVersionPackageAnalyzeResult) GetStruct() AnalyzeResult {
	return r
}

func (r MultiVersionPackageAnalyzeResult) OutputText(analyzeType string) error {
	analysis := r.Analysis
	strAnalysis := stringifyMultiVersionPackages(analysis)

	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    map[string]map[string]StrPackageInfo
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    strAnalysis,
	}
	return TemplateOutput(strResult, "MultiVersionPackageAnalyze")
}

func stringifyMultiVersionPackages(packages map[string]map[string]PackageInfo) map[string]map[string]StrPackageInfo {
	strPackages := map[string]map[string]StrPackageInfo{}
	for pack, versionMap := range packages {
		strPackages[pack] = stringifyPackages(versionMap)
	}
	return strPackages
}

type SingleVersionPackageAnalyzeResult struct {
	Image       string
	AnalyzeType string
	Analysis    map[string]PackageInfo
}

func (r SingleVersionPackageAnalyzeResult) GetStruct() AnalyzeResult {
	return r
}

func (r SingleVersionPackageAnalyzeResult) OutputText(diffType string) error {
	analysis := r.Analysis
	strAnalysis := stringifyPackages(analysis)

	strResult := struct {
		Image       string
		AnalyzeType string
		Analysis    map[string]StrPackageInfo
	}{
		Image:       r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis:    strAnalysis,
	}
	return TemplateOutput(strResult, "SingleVersionPackageAnalyze")
}

type StrPackageInfo struct {
	Version string
	Size    string
}

func stringifyPackageInfo(info PackageInfo) StrPackageInfo {
	return StrPackageInfo{Version: info.Version, Size: stringifySize(info.Size)}
}

func stringifySize(size int64) string {
	strSize := "unknown"
	if size != -1 {
		strSize = bytefmt.ByteSize(uint64(size))
	}
	return strSize
}

func stringifyPackages(packages map[string]PackageInfo) map[string]StrPackageInfo {
	strPackages := map[string]StrPackageInfo{}
	for pack, info := range packages {
		strInfo := stringifyPackageInfo(info)
		strPackages[pack] = strInfo
	}
	return strPackages
}

type FileAnalyzeResult struct {
	Image       string
	AnalyzeType string
	Analysis    []DirectoryEntry
}

func (r FileAnalyzeResult) GetStruct() AnalyzeResult {
	return r
}

func (r FileAnalyzeResult) OutputText(analyzeType string) error {
	analysis := r.Analysis
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
