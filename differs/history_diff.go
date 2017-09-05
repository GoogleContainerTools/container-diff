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
	"strings"

	"github.com/GoogleCloudPlatform/container-diff/utils"
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

func (a HistoryAnalyzer) Diff(image1, image2 utils.Image) (utils.Result, error) {
	diff, err := getHistoryDiff(image1, image2)
	return &utils.HistDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: "History",
		Diff:     diff,
	}, err
}

func (a HistoryAnalyzer) Analyze(image utils.Image) (utils.Result, error) {
	history := getHistoryList(image.Config.History)
	result := utils.ListAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: "History",
		Analysis:    history,
	}
	return &result, nil
}

func getHistoryDiff(image1, image2 utils.Image) (HistDiff, error) {
	history1 := getHistoryList(image1.Config.History)
	history2 := getHistoryList(image2.Config.History)

	adds := utils.GetAdditions(history1, history2)
	dels := utils.GetDeletions(history1, history2)
	diff := HistDiff{adds, dels}
	return diff, nil
}

func getHistoryList(historyItems []utils.ImageHistoryItem) []string {
	strhistory := make([]string, len(historyItems))
	for i, layer := range historyItems {
		strhistory[i] = strings.TrimSpace(layer.CreatedBy)
	}
	return strhistory
}
