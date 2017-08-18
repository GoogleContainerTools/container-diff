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
	return TemplateOutput(r, "MultiVersionPackageAnalyze")
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
	return TemplateOutput(r, "SingleVersionPackageAnalyze")
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
		Image string
		AnalyzeType string
		Analysis []StrDirectoryEntry
	}{
		Image: r.Image,
		AnalyzeType: r.AnalyzeType,
		Analysis: strAnalysis,
	}
	return TemplateOutput(strResult, "FileAnalyze")
}

type StrDirectoryEntry struct {
	Name string
	Size string
}

func stringifyDirectoryEntries(entries []DirectoryEntry) (strEntries []StrDirectoryEntry) {
	for _, entry := range entries {
		size := entry.Size
		strSize := "unknown"
		if size != -1 {
			strSize = bytefmt.ByteSize(uint64(size))
		}

		strEntry := StrDirectoryEntry{Name: entry.Name, Size: strSize}
		strEntries = append(strEntries, strEntry)
	}
	return
}
