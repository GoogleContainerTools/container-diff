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
	"regexp"
	"strings"

	"github.com/containers/image/docker/daemon"
	"github.com/containers/image/docker/reference"
	"github.com/golang/glog"
)

const DaemonPrefix = "daemon://"

type DaemonPrepper struct {
	*ImagePrepper
}

func (p DaemonPrepper) Name() string {
	return "Local Daemon"
}

func (p DaemonPrepper) GetSource() string {
	return p.ImagePrepper.Source
}

func (p DaemonPrepper) SupportsImage() bool {
	// will fail on strings prefixed with 'remote://'
	remoteRegex := regexp.MustCompile(RemotePrefix + ".*")
	if match := remoteRegex.MatchString(p.ImagePrepper.Source); match || IsTar(p.ImagePrepper.Source) {
		return false
	}
	strippedSource := strings.Replace(p.ImagePrepper.Source, DaemonPrefix, "", -1)
	_, err := reference.Parse(strippedSource)
	if err == nil {
		// strip prefix off image source for later use
		p.ImagePrepper.Source = strippedSource
		return true
	}
	return false
}

func (p DaemonPrepper) GetFileSystem() (string, error) {
	ref, err := daemon.ParseReference(p.Source)
	if err != nil {
		return "", err
	}

	return getFileSystemFromReference(ref, p.Source)
}

func (p DaemonPrepper) GetConfig() (ConfigSchema, error) {
	inspect, _, err := p.Client.ImageInspectWithRaw(context.Background(), p.Source)
	if err != nil {
		return ConfigSchema{}, err
	}

	config := ConfigObject{
		Env: inspect.Config.Env,
	}
	history := p.GetHistory()
	return ConfigSchema{
		Config:  config,
		History: history,
	}, nil
}

func (p DaemonPrepper) GetHistory() []ImageHistoryItem {
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
