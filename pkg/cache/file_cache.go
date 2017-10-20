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

package cache

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type FileCache struct {
	RootDir string
}

func NewFileCache(dir string) (*FileCache, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &FileCache{RootDir: dir}, nil
}

func (c *FileCache) HasLayer(layerId string) bool {
	_, err := os.Stat(filepath.Join(c.RootDir, layerId))
	return !os.IsNotExist(err)
}

func (c *FileCache) SetLayer(layerId string, r io.Reader) (io.ReadCloser, error) {
	fullpath := filepath.Join(c.RootDir, layerId)
	entry, err := os.Create(fullpath)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(entry, r); err != nil {
		return nil, err
	}
	return c.GetLayer(layerId)
}

func (c *FileCache) GetLayer(layerId string) (io.ReadCloser, error) {
	logrus.Infof("retrieving layer %s from cache", layerId)
	return os.Open(filepath.Join(c.RootDir, layerId))
}

func (c *FileCache) Invalidate(layerId string) error {
	return os.RemoveAll(filepath.Join(c.RootDir, layerId))
}
