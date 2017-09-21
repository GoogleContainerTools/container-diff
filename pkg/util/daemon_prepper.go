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
	"strings"

	"github.com/containers/image/docker/daemon"
	"github.com/containers/image/docker/reference"
	"github.com/golang/glog"
)

const DaemonPrefix = "daemon://"

type DaemonPrepper struct {
	ImagePrepper
}

func (p DaemonPrepper) Name() string {
	return "Local Daemon"
}

func (p DaemonPrepper) RawSource() string {
	return strings.Replace(p.Source, DaemonPrefix, "", -1)
}

func (p DaemonPrepper) SupportsImage() bool {
	// will fail on strings prefixed with 'remote://'
	if strings.HasPrefix(p.Source, RemotePrefix) || IsTar(p.Source) {
		return false
	}
	_, err := reference.Parse(p.RawSource())
	return (err == nil)
}

func (p DaemonPrepper) GetFileSystem() (string, error) {
	ref, err := daemon.ParseReference(p.RawSource())
	if err != nil {
		return "", err
	}

	return getFileSystemFromReference(ref, p.RawSource())
}

func (p DaemonPrepper) GetConfig() (ConfigSchema, error) {
	ref, err := daemon.ParseReference(p.RawSource())
	if err != nil {
		return ConfigSchema{}, err
	}

	return getConfigFromReference(ref, p.RawSource())
}

func (p DaemonPrepper) GetHistory() []ImageHistoryItem {
	history, err := p.Client.ImageHistory(context.Background(), p.RawSource())
	if err != nil {
		glog.Error("Could not obtain image history for %s: %s", p.RawSource(), err)
	}
	historyItems := []ImageHistoryItem{}
	for _, item := range history {
		historyItems = append(historyItems, ImageHistoryItem{CreatedBy: item.CreatedBy})
	}
	return historyItems
}
