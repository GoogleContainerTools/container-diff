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
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/containers/image/types"
	digest "github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
)

type FileCache struct {
	RootDir string
	Ref     types.ImageReference
	src     types.ImageSource
}

func NewFileCache(dir string, ref types.ImageReference, src types.ImageSource) (*FileCache, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	return &FileCache{
		RootDir: dir,
		Ref:     ref,
		src:     src,
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
	entry, err := os.Create(fullpath)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(entry, r); err != nil {
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

// Implement types.ImageSource

func (c *FileCache) Reference() types.ImageReference {
	return c.Ref
}

func (c *FileCache) Close() error {
	return nil
}

func (c *FileCache) GetManifest() ([]byte, string, error) {
	return c.src.GetManifest()
}

func (c *FileCache) GetTargetManifest(digest digest.Digest) ([]byte, string, error) {
	return c.GetTargetManifest(digest)
}

// GetBlob returns a stream for the specified blob, and the blob’s size (or -1 if unknown).
// The Digest field in BlobInfo is guaranteed to be provided; Size may be -1.
func (c *FileCache) GetBlob(bi types.BlobInfo) (io.ReadCloser, int64, error) {
	if c.HasLayer(bi) {
		r, err := c.GetLayer(bi)
		return r, bi.Size, err
	}
	// Add to the cache then return
	r, size, err := c.src.GetBlob(bi)
	if err != nil {
		return nil, 0, err
	}
	r, err = c.SetLayer(bi, r)
	if err != nil {
		return nil, 0, c.Invalidate(bi)
	}
	return r, size, err
}

// GetSignatures returns the image's signatures.  It may use a remote (= slow) service.
func (c *FileCache) GetSignatures(ctx context.Context) ([][]byte, error) {
	return c.src.GetSignatures(ctx)
}
