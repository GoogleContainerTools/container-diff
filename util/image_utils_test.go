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
	"testing"

	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
)

type imageTestPair struct {
	input          string
	expectedOutput bool
}

func TestCheckImageID(t *testing.T) {
	for _, test := range []imageTestPair{
		{input: "gcr.io/repo/image", expectedOutput: false},
		{input: "testTars/la-croix1.tar", expectedOutput: false},
	} {
		prepper := pkgutil.DaemonPrepper{
			&pkgutil.ImagePrepper{
				Source: test.input,
			},
		}
		output := prepper.SupportsImage()
		if output != test.expectedOutput {
			if test.expectedOutput {
				t.Errorf("Expected input to be image ID but %s tested false", test.input)
			} else {
				t.Errorf("Didn't expect input to be an image ID but %s tested true", test.input)
			}
		}
	}
}

func TestCheckImageTar(t *testing.T) {
	for _, test := range []imageTestPair{
		{input: "gcr.io/repo/image", expectedOutput: false},
		{input: "testTars/la-croix1.tar", expectedOutput: true},
	} {
		output := pkgutil.CheckTar(test.input)
		if output != test.expectedOutput {
			if test.expectedOutput {
				t.Errorf("Expected input to be a tar file but %s tested false", test.input)
			} else {
				t.Errorf("Didn't expect input to be a tar file but %s tested true", test.input)
			}
		}
	}
}

func TestCheckImageURL(t *testing.T) {
	for _, test := range []imageTestPair{
		{input: "gcr.io/repo/image", expectedOutput: true},
		{input: "testTars/la-croix1.tar", expectedOutput: false},
	} {
		prepper := pkgutil.CloudPrepper{
			&pkgutil.ImagePrepper{
				Source: test.input,
			},
		}
		output := prepper.SupportsImage()
		if output != test.expectedOutput {
			if test.expectedOutput {
				t.Errorf("Expected input to be a tar file but %s tested false", test.input)
			} else {
				t.Errorf("Didn't expect input to be a tar file but %s tested true", test.input)
			}
		}
	}
}
