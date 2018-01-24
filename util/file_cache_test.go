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
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/opencontainers/go-digest"

	"github.com/GoogleCloudPlatform/container-diff/pkg/cache"
	"github.com/containers/image/types"
)

func cacheAndTest(c *cache.FileCache, t *testing.T, testStr string, layerId string) {
	byteArr := []byte(testStr)
	r := bytes.NewReader(byteArr)

	bi := types.BlobInfo{
		Digest: digest.Digest(layerId),
	}
	if c.HasLayer(bi) {
		t.Errorf("cache already has test layer %s", layerId)
	}
	c.SetLayer(bi, r)

	if !c.HasLayer(bi) {
		t.Errorf("layer %s not successfully cached", layerId)
	}
	cachedLayer, err := c.GetLayer(bi)
	if err != nil {
		t.Errorf(err.Error())
	}
	cachedData, err := ioutil.ReadAll(cachedLayer)
	cachedStr := string(cachedData)
	if cachedStr != testStr {
		t.Errorf("cached data %s does not match original: %s", cachedStr, testStr)
	}
}

func TestCache(t *testing.T) {
	cacheDir, err := ioutil.TempDir("", ".cache")
	if err != nil {
		t.Fatalf("error when creating cache directory: %s", err.Error())
	}
	defer os.RemoveAll(cacheDir)
	c := &cache.FileCache{RootDir: cacheDir}
	testRuns := []struct {
		Name    string
		Data    string
		LayerId string
	}{
		{"real data", "this is a test of caching some bytes. this could be any data.", "sha256:realdata"},
		{"empty data", "", "sha256:emptydata"},
	}
	for _, test := range testRuns {
		t.Run(test.Name, func(t *testing.T) {
			cacheAndTest(c, t, test.Data, test.LayerId)
		})
	}
}
