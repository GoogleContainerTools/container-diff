package utils

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
	return TemplateOutput(r)
}

type MultiPackageAnalyzeResult struct {
	Image       string
	AnalyzeType string
	Analysis    map[string]map[string]PackageInfo
}

func (r MultiPackageAnalyzeResult) GetStruct() AnalyzeResult {
	return r
}

func (r MultiPackageAnalyzeResult) OutputText(analyzeType string) error {
	return TemplateOutput(r)
}

type PackageAnalyzeResult struct {
	Image       string
	AnalyzeType string
	Analysis    map[string]PackageInfo
}

func (r PackageAnalyzeResult) GetStruct() AnalyzeResult {
	return r
}

func (r PackageAnalyzeResult) OutputText(diffType string) error {
	return TemplateOutput(r)
}
