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
	return TemplateOutput(r)
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
