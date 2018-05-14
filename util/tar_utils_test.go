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

package util

import (
	"testing"

	pkgutil "github.com/GoogleContainerTools/container-diff/pkg/util"
)

func TestIsTar(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{input: "/testTar/la-croix1.tar", expected: true},
		{input: "/testTar/la-croix1-actual", expected: false},
	}
	for _, test := range testCases {
		actual := pkgutil.IsTar(test.input)
		if test.expected != actual {
			t.Errorf("Expected: %t but got: %t", test.expected, actual)
		}
	}
}
