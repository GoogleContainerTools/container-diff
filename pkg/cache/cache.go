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
)

type FileCache struct {
	RootDir string
}

func NewFileCache(rootDir string) (*FileCache, error) {
	if err := os.MkdirAll(rootDir, 0700); err != nil {
		return nil, err
	}
	return &FileCache{RootDir: rootDir}, nil
}

func (f *FileCache) HasLayer(layerId string) bool {
	_, err := os.Stat(filepath.Join(f.RootDir, layerId))
	return !os.IsNotExist(err)
}

func (f *FileCache) SetLayer(layerId string, r io.Reader) (io.ReadCloser, error) {
	path := filepath.Join(f.RootDir, layerId)
	l, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(l, r); err != nil {
		return nil, err
	}

	return f.GetLayer(layerId)
}

func (f *FileCache) GetLayer(layerId string) (io.ReadCloser, error) {
	path := filepath.Join(f.RootDir, layerId)
	return os.Open(path)
}
