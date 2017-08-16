package differs

import (
	"sort"

	"github.com/GoogleCloudPlatform/container-diff/utils"
)

type FileDiffer struct {
}

// FileDiff diffs two packages and compares their contents
func (d FileDiffer) Diff(image1, image2 utils.Image) (utils.DiffResult, error) {
	diff, err := diffImageFiles(image1, image2)
	return &utils.DirDiffResult{DiffType: "FileDiffer", Diff: diff}, err
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
		Image1: image1.Source,
		Image2: image2.Source,
		Adds:   adds,
		Dels:   dels,
		Mods:   []string{},
	}
	return diff, nil
}
