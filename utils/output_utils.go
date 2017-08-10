package utils

type DiffResult interface {
	GetStruct() DiffResult
	OutputText(diffType string) error
}

type MultiVersionPackageDiffResult struct {
	DiffType string
	Diff     MultiVersionPackageDiff
}

func (m MultiVersionPackageDiffResult) GetStruct() DiffResult {
	return m
}

func (m MultiVersionPackageDiffResult) OutputText(diffType string) error {
	return TemplateOutput(m)
}

type PackageDiffResult struct {
	DiffType string
	Diff     PackageDiff
}

func (m PackageDiffResult) GetStruct() DiffResult {
	return m
}

func (m PackageDiffResult) OutputText(diffType string) error {
	return TemplateOutput(m)
}

type HistDiffResult struct {
	DiffType string
	Diff     HistDiff
}

func (m HistDiffResult) GetStruct() DiffResult {
	return m
}

func (m HistDiffResult) OutputText(diffType string) error {
	return TemplateOutput(m)
}

type DirDiffResult struct {
	DiffType string
	Diff     DirDiff
}

func (m DirDiffResult) GetStruct() DiffResult {
	return m
}

func (m DirDiffResult) OutputText(diffType string) error {
	return TemplateOutput(m)
}
