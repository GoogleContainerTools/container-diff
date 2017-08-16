package differs

import (
	"sort"

	"github.com/GoogleCloudPlatform/container-diff/utils"
)

type FileAnalyzer struct {
}

// FileDiff diffs two packages and compares their contents
func (a FileAnalyzer) Diff(image1, image2 utils.Image) (utils.DiffResult, error) {
	diff, err := diffImageFiles(image1, image2)
	return &utils.DirDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: "File",
		Diff:     diff,
	}, err
}

func (a FileAnalyzer) Analyze(image utils.Image) (utils.AnalyzeResult, error) {
	var result utils.ListAnalyzeResult

	imgDir, err := utils.GetDirectory(image.FSPath, true)
	if err != nil {
		return result, err
	}

	result.Image = image.Source
	result.AnalyzeType = "File"
	result.Analysis = imgDir.Content
	return &result, err
}

func diffImageFiles(image1, image2 utils.Image) (utils.DirDiff, error) {
	img1 := image1.FSPath
	img2 := image2.FSPath

	var diff utils.DirDiff

	img1Dir, err := utils.GetDirectory(img1, true)
	if err != nil {
		return diff, err
	}
	img2Dir, err := utils.GetDirectory(img2, true)
	if err != nil {
		return diff, err
	}

	adds := utils.GetAddedEntries(img1Dir, img2Dir)
	sort.Strings(adds)
	dels := utils.GetDeletedEntries(img1Dir, img2Dir)
	sort.Strings(dels)

	diff = utils.DirDiff{
		Adds: adds,
		Dels: dels,
		Mods: []string{},
	}
	return diff, nil
}
