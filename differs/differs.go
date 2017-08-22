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
	DiffTypes []Analyzer
}

type SingleRequest struct {
	Image        utils.Image
	AnalyzeTypes []Analyzer
}

type Analyzer interface {
	Diff(image1, image2 utils.Image) (utils.Result, error)
	Analyze(image utils.Image) (utils.Result, error)
}

var analyzers = map[string]Analyzer{
	"history": HistoryAnalyzer{},
	"file":    FileAnalyzer{},
	"apt":     AptAnalyzer{},
	"pip":     PipAnalyzer{},
	"node":    NodeAnalyzer{},
}

func (req DiffRequest) GetDiff() (map[string]utils.Result, error) {
	img1 := req.Image1
	img2 := req.Image2
	diffs := req.DiffTypes

	results := map[string]utils.Result{}
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

func (req SingleRequest) GetAnalysis() (map[string]utils.Result, error) {
	img := req.Image
	analyses := req.AnalyzeTypes

	results := map[string]utils.Result{}
	for _, analyzer := range analyses {
		analyzeName := reflect.TypeOf(analyzer).Name()
		if analysis, err := analyzer.Analyze(img); err == nil {
			results[analyzeName] = analysis
		} else {
			glog.Errorf("Error getting analysis with %s: %s", analyzeName, err)
		}
	}

	var err error
	if len(results) == 0 {
		err = fmt.Errorf("Could not perform analysis on %s", img)
	} else {
		err = nil
	}

	return results, err
}

func GetAnalyzers(analyzeNames []string) (analyzeFuncs []Analyzer, err error) {
	for _, name := range analyzeNames {
		if a, exists := analyzers[name]; exists {
			analyzeFuncs = append(analyzeFuncs, a)
		} else {
			glog.Errorf("Unknown analyzer/differ specified", name)
		}
	}
	if len(analyzeFuncs) == 0 {
		err = errors.New("No known analyzers/differs specified")
	}
	return
}
