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

package util

import (
	"context"

	"github.com/containers/image/docker/daemon"
	"github.com/golang/glog"
)

type DaemonPrepper struct {
	ImagePrepper
}

func (p DaemonPrepper) getFileSystem() (string, error) {
	ref, err := daemon.ParseReference(p.Source)
	if err != nil {
		return "", err
	}

	return getFileSystemFromReference(ref, p.Source)
}

func (p DaemonPrepper) getConfig() (ConfigSchema, error) {
	ref, err := daemon.ParseReference(p.Source)
	if err != nil {
		return ConfigSchema{}, err
	}

	return getConfigFromReference(ref, p.Source)
}

func (p DaemonPrepper) getHistory() []ImageHistoryItem {
	history, err := p.Client.ImageHistory(context.Background(), p.Source)
	if err != nil {
		glog.Error("Could not obtain image history for %s: %s", p.Source, err)
	}
	historyItems := []ImageHistoryItem{}
	for _, item := range history {
		historyItems = append(historyItems, ImageHistoryItem{CreatedBy: item.CreatedBy})
	}
	return historyItems
}