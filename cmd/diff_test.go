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
	"testing"
)

var diffArgNumTests = []testpair{
	{[]string{}, true},
	{[]string{"one"}, true},
	{[]string{"one", "two"}, false},
	{[]string{"one", "two", "three"}, true},
}

func TestDiffArgNum(t *testing.T) {
	for _, test := range diffArgNumTests {
		err := checkDiffArgNum(test.input)
		checkError(t, err, test.shouldError)
	}
}

type imageDiff struct {
	image1      string
	image2      string
	shouldError bool
}

var imageDiffs = []imageDiff{
	{"", "", true},
	{"gcr.io/google-appengine/python", "gcr.io/google-appengine/debian9", false},
	{"gcr.io/google-appengine/python", "cats", true},
	{"mcr.microsoft.com/mcr/hello-world:latest", "mcr.microsoft.com/mcr/hello-world:latest", false},
}

func TestDiffImages(t *testing.T) {
	for _, test := range imageDiffs {
		err := diffImages(test.image1, test.image2, []string{"apt"})
		checkError(t, err, test.shouldError)
		err = diffImages(test.image1, test.image2, []string{"metadata"})
		checkError(t, err, test.shouldError)
	}
}

func checkError(t *testing.T, err error, shouldError bool) {
	if (err == nil) == shouldError {
		if shouldError {
			t.Errorf("expected error but got none")
		} else {
			t.Errorf("got unexpected error: %s", err)
		}
	}
}
