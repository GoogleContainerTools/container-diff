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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/GoogleCloudPlatform/container-diff/pkg/image"
	"github.com/containers/image/types"
	"github.com/sirupsen/logrus"
)

type FileCache struct {
	RootDir string
	*image.ProxySource
}

func NewFileCache(dir string, ref types.ImageReference) (*FileCache, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	ps, err := image.NewProxySource(ref)
	if err != nil {
		return nil, err
	}

	return &FileCache{
		RootDir:     dir,
		ProxySource: ps,
	}, nil
}

func (c *FileCache) HasLayer(layer types.BlobInfo) bool {
	layerId := layer.Digest.String()
	_, err := os.Stat(filepath.Join(c.RootDir, layerId))
	return !os.IsNotExist(err)
}

func (c *FileCache) SetLayer(layer types.BlobInfo, r io.Reader) (io.ReadCloser, error) {
	layerId := layer.Digest.String()
	fullpath := filepath.Join(c.RootDir, layerId)
	// Write the entry atomically. First write it to a temporary name, then rename to the correct one.
	f, err := ioutil.TempFile(c.RootDir, "")
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(f, r); err != nil {
		return nil, err
	}
	if err := os.Rename(f.Name(), fullpath); err != nil {
		return nil, err
	}
	return c.GetLayer(layer)
}

func (c *FileCache) GetLayer(layer types.BlobInfo) (io.ReadCloser, error) {
	layerId := layer.Digest.String()
	logrus.Infof("retrieving layer %s from cache", layerId)
	return os.Open(filepath.Join(c.RootDir, layerId))
}

func (c *FileCache) Invalidate(layer types.BlobInfo) error {
	layerId := layer.Digest.String()
	return os.RemoveAll(filepath.Join(c.RootDir, layerId))
}

// GetBlob returns a stream for the specified blob, and the blob’s size (or -1 if unknown).
// The Digest field in BlobInfo is guaranteed to be provided; Size may be -1.
func (c *FileCache) GetBlob(bi types.BlobInfo) (io.ReadCloser, int64, error) {
	if c.HasLayer(bi) {
		r, err := c.GetLayer(bi)
		return r, bi.Size, err
	}
	// Add to the cache then return
	r, size, err := c.ProxySource.GetBlob(bi)
	if err != nil {
		return nil, 0, err
	}
	r, err = c.SetLayer(bi, r)
	if err != nil {
		logrus.Errorf("Error setting layer %s in cache: %v", bi.Digest, err)
		return nil, 0, c.Invalidate(bi)
	}
	return r, size, err
}
