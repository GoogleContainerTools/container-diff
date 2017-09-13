/*
Copyright 2017 Google, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package differs

import (
	"errors"
	"fmt"

	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/GoogleCloudPlatform/container-diff/util"
	"github.com/golang/glog"
)

type DiffRequest struct {
	Image1    pkgutil.Image
	Image2    pkgutil.Image
	DiffTypes []Analyzer
}

type SingleRequest struct {
	Image        pkgutil.Image
	AnalyzeTypes []Analyzer
}

type Analyzer interface {
	Diff(image1, image2 pkgutil.Image) (util.Result, error)
	Analyze(image pkgutil.Image) (util.Result, error)
	Name() string
}

var Analyzers = map[string]Analyzer{
	"history": HistoryAnalyzer{},
	"file":    FileAnalyzer{},
	"apt":     AptAnalyzer{},
	"pip":     PipAnalyzer{},
	"node":    NodeAnalyzer{},
}

func (req DiffRequest) GetDiff() (map[string]util.Result, error) {
	img1 := req.Image1
	img2 := req.Image2
	diffs := req.DiffTypes

	results := map[string]util.Result{}
	for _, differ := range diffs {
		if diff, err := differ.Diff(img1, img2); err == nil {
			results[differ.Name()] = diff
		} else {
			glog.Errorf("Error getting diff with %s: %s", differ.Name(), err)
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

func (req SingleRequest) GetAnalysis() (map[string]util.Result, error) {
	img := req.Image
	analyses := req.AnalyzeTypes

	results := map[string]util.Result{}
	for _, analyzer := range analyses {
		analyzeName := analyzer.Name()
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
		if a, exists := Analyzers[name]; exists {
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
