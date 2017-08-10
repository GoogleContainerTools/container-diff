package differs

import (
	"strings"

	"github.com/GoogleCloudPlatform/container-diff/utils"
)

type HistoryDiffer struct {
}

func (d HistoryDiffer) Diff(image1, image2 utils.Image) (utils.DiffResult, error) {
	diff, err := getHistoryDiff(image1, image2)
	return &utils.HistDiffResult{DiffType: "HistoryDiffer", Diff: diff}, err
}

func getHistoryDiff(image1, image2 utils.Image) (utils.HistDiff, error) {
	history1 := getHistoryList(image1.Config.History)
	history2 := getHistoryList(image2.Config.History)

	adds := utils.GetAdditions(history1, history2)
	dels := utils.GetDeletions(history1, history2)
	diff := utils.HistDiff{image1.Source, image2.Source, adds, dels}
	return diff, nil
}

func getHistoryList(historyItems []utils.ImageHistoryItem) []string {
	strhistory := make([]string, len(historyItems))
	for i, layer := range historyItems {
		strhistory[i] = strings.TrimSpace(layer.CreatedBy)
	}
	return strhistory
}
