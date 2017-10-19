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
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/container-diff/pkg/cache"
)

var c *cache.FileCache

func TestCacheData(t *testing.T) {
	data := "this is a test of caching some bytes. this could be any data."
	cacheAndTest(t, data, "sha256:realdata")
}

func TestCacheEmpty(t *testing.T) {
	cacheAndTest(t, "", "sha256:emptydata")
}

func cacheAndTest(t *testing.T, testStr string, layerId string) {
	byteArr := []byte(testStr)
	r := bytes.NewReader(byteArr)

	if c.HasLayer(layerId) {
		t.Errorf("cache already has test layer %s", layerId)
	}
	c.SetLayer(layerId, r)

	if !c.HasLayer(layerId) {
		t.Errorf("layer %s not successfully cached", layerId)
	}
	cachedLayer, err := c.GetLayer(layerId)
	if err != nil {
		t.Errorf(err.Error())
	}
	cachedData, err := ioutil.ReadAll(cachedLayer)
	cachedStr := string(cachedData)
	if cachedStr != testStr {
		t.Errorf("cached data %s does not match original: %s", cachedStr, testStr)
	}
}

func TestMain(m *testing.M) {
	var err error
	cacheDir, err := ioutil.TempDir("", ".cache")
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	defer os.RemoveAll(cacheDir)
	c, err = cache.NewFileCache(cacheDir)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	os.Exit(m.Run())
}
