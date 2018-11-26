/*
Copyright 2018 Google, Inc. All rights reserved.

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

package cmd

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
)

type testpair struct {
	input       []string
	shouldError bool
}

func TestCacheDir(t *testing.T) {
	homeDir, err := homedir.Dir()
	if err != nil {
		t.Errorf("error getting home dir: %s", err.Error())
	}
	tests := []struct {
		name        string
		cliFlag     string
		envVar      string
		expectedDir string
		imageName   string
	}{
		{
			name:        "default cache is at $HOME",
			cliFlag:     "",
			envVar:      "",
			expectedDir: filepath.Join(homeDir, ".container-diff", "cache"),
			imageName:   "pancakes",
		},
		{
			name:        "setting cache via --cache-dir",
			cliFlag:     "/tmp",
			envVar:      "",
			expectedDir: "/tmp/.container-diff/cache",
			imageName:   "pancakes",
		},
		{
			name:        "setting cache via CONTAINER_DIFF_CACHEDIR",
			cliFlag:     "",
			envVar:      "/tmp",
			expectedDir: "/tmp/.container-diff/cache",
			imageName:   "pancakes",
		},
		{
			name:        "command line --cache-dir takes preference to CONTAINER_DIFF_CACHEDIR",
			cliFlag:     "/tmp",
			envVar:      "/opt",
			expectedDir: "/tmp/.container-diff/cache",
			imageName:   "pancakes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set any environment variables
			if tt.envVar != "" {
				os.Setenv("CONTAINER_DIFF_CACHEDIR", tt.envVar)
			}
			// Set global flag for cache based on --cache-dir
			cacheDir = tt.cliFlag

			// call getCacheDir and make sure return is equal to expected
			actualDir, err := getCacheDir(tt.imageName)
			if err != nil {
				t.Errorf("Error getting cache dir %s: %s", tt.name, err.Error())
			}

			if path.Dir(actualDir) != tt.expectedDir {
				t.Errorf("%s\nexpected: %v\ngot: %v", tt.name, tt.expectedDir, actualDir)
			}
		},
		)
	}
}
