package differs

import (
	"os"
	"sort"

	"github.com/GoogleCloudPlatform/container-diff/utils"
	"github.com/golang/glog"
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

	target1 := "j1.json"
	err := utils.DirToJSON(img1, target1, true)
	if err != nil {
		return diff, err
	}
	target2 := "j2.json"
	err = utils.DirToJSON(img2, target2, true)
	if err != nil {
		return diff, err
	}
	img1Dir, err := utils.GetDirectory(target1)
	if err != nil {
		return diff, err
	}
	img2Dir, err := utils.GetDirectory(target2)
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
	}

	err = os.Remove(target1)
	if err != nil {
		glog.Error(err)
	}
	err = os.Remove(target2)
	if err != nil {
		glog.Error(err)
	}
	return diff, nil
}
