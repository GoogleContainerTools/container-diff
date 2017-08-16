package utils

type DiffResult interface {
	GetStruct() DiffResult
	OutputText(diffType string) error
}

type MultiPackageDiffResult struct {
	Image1   string
	Image2   string
	DiffType string
	Diff     MultiVersionPackageDiff
}

func (r MultiPackageDiffResult) GetStruct() DiffResult {
	return r
}

func (r MultiPackageDiffResult) OutputText(diffType string) error {
	return TemplateOutput(r)
}

type PackageDiffResult struct {
	Image1   string
	Image2   string
	DiffType string
	Diff     PackageDiff
}

func (r PackageDiffResult) GetStruct() DiffResult {
	return r
}

func (r PackageDiffResult) OutputText(diffType string) error {
	return TemplateOutput(r)
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
	return TemplateOutput(r)
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
	return TemplateOutput(r)
}
