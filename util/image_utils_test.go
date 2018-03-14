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

	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
)

func TestImageTags(t *testing.T) {
	tests := []struct {
		image  string
		hasTag bool
	}{
		{
			image:  "gcr.io/test_image/foo:latest",
			hasTag: true,
		},
		{
			image:  "gcr.io/test_image/foo:",
			hasTag: false,
		},
		{
			image:  "daemon://gcr.io/test_image/foo:test",
			hasTag: true,
		},
		{
			image:  "remote://gcr.io/test_image_foo",
			hasTag: false,
		},
	}

	for _, test := range tests {
		if pkgutil.HasTag(test.image) != test.hasTag {
			t.Errorf("Error checking tag on image %s", test.image)
		}
	}
}
