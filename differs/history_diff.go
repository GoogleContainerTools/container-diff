/*
Copyright 2018 Google, Inc. All rights reserved.

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
	"strings"

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
	"github.com/GoogleContainerTools/container-diff/util"
	"github.com/google/go-containerregistry/v1"
)

type HistoryAnalyzer struct {
}

type HistDiff struct {
	Adds []string
	Dels []string
}

func (a HistoryAnalyzer) Name() string {
	return "HistoryAnalyzer"
}

func (a HistoryAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff, err := getHistoryDiff(image1, image2)
	return &util.HistDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: "History",
		Diff:     diff,
	}, err
}

func (a HistoryAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	c, err := image.Image.ConfigFile()
	if err != nil {
		return util.ListAnalyzeResult{}, err
	}
	history := getHistoryList(c.History)
	result := util.ListAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: "History",
		Analysis:    history,
	}
	return &result, nil
}

func getHistoryDiff(image1, image2 pkgutil.Image) (HistDiff, error) {
	c1, err := image1.Image.ConfigFile()
	if err != nil {
		return HistDiff{}, err
	}
	c2, err := image2.Image.ConfigFile()
	if err != nil {
		return HistDiff{}, err
	}
	history1 := getHistoryList(c1.History)
	history2 := getHistoryList(c2.History)

	adds := util.GetAdditions(history1, history2)
	dels := util.GetDeletions(history1, history2)
	diff := HistDiff{adds, dels}
	return diff, nil
}

func getHistoryList(historyItems []v1.History) []string {
	strhistory := make([]string, len(historyItems))
	for i, layer := range historyItems {
		strhistory[i] = strings.TrimSpace(layer.CreatedBy)
	}
	return strhistory
}
