package differs

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/GoogleCloudPlatform/container-diff/utils"
	"github.com/golang/glog"
)

type DiffRequest struct {
	Image1    utils.Image
	Image2    utils.Image
	DiffTypes []Differ
}

type Differ interface {
	Diff(image1, image2 utils.Image) (utils.DiffResult, error)
}

var diffs = map[string]Differ{
	"history": HistoryDiffer{},
	"file":    FileDiffer{},
	"apt":     AptDiffer{},
	"pip":     PipDiffer{},
	"node":    NodeDiffer{},
}

func (diff DiffRequest) GetDiff() (map[string]utils.DiffResult, error) {
	img1 := diff.Image1
	img2 := diff.Image2
	diffs := diff.DiffTypes

	results := map[string]utils.DiffResult{}
	for _, differ := range diffs {
		differName := reflect.TypeOf(differ).Name()
		if diff, err := differ.Diff(img1, img2); err == nil {
			results[differName] = diff
		} else {
			glog.Errorf("Error getting diff with %s: %s", differName, err)
		}
	}

	var err error
	if len(results) == 0 {
		err = fmt.Errorf("Could not perform diff on %s and %s", img1, img2)
	} else {
		err = nil
	}

	return results, err
}

func GetDiffers(diffNames []string) (diffFuncs []Differ, err error) {
	for _, diffName := range diffNames {
		if d, exists := diffs[diffName]; exists {
			diffFuncs = append(diffFuncs, d)
		} else {
			glog.Errorf("Unknown differ specified", diffName)
		}
	}
	if len(diffFuncs) == 0 {
		err = errors.New("No known differs specified")
	}
	return
}
