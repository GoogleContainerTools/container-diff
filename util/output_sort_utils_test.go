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
	"reflect"
	"testing"

	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
)

var packageTests = [][]PackageOutput{
	{
		{Name: "a", Version: "1.2", Size: 10},
		{Name: "b", Version: "1.5", Size: 12},
		{Name: "c", Version: "1.4", Size: 20},
	},
	{
		{Name: "a", Version: "1.2", Size: 10},
		{Name: "b", Version: "1.5", Size: 12},
		{Name: "c", Version: "1.4", Size: 12},
	},
	{
		{Name: "a", Version: "1.2", Size: 10},
		{Name: "a", Version: "1.4", Size: 20},
		{Name: "a", Version: "1.2", Size: 15},
	},
}

func TestSortPackageOutput(t *testing.T) {
	for _, test := range []struct {
		input    []PackageOutput
		sortBy   func(a, b *PackageOutput) bool
		expected []PackageOutput
	}{
		{
			input:  packageTests[0],
			sortBy: packageSizeSort,
			expected: []PackageOutput{
				{Name: "c", Version: "1.4", Size: 20},
				{Name: "b", Version: "1.5", Size: 12},
				{Name: "a", Version: "1.2", Size: 10},
			},
		},
		{
			input:  packageTests[0],
			sortBy: packageNameSort,
			expected: []PackageOutput{
				{Name: "a", Version: "1.2", Size: 10},
				{Name: "b", Version: "1.5", Size: 12},
				{Name: "c", Version: "1.4", Size: 20},
			},
		},
		{
			input:  packageTests[1],
			sortBy: packageSizeSort,
			expected: []PackageOutput{
				{Name: "b", Version: "1.5", Size: 12},
				{Name: "c", Version: "1.4", Size: 12},
				{Name: "a", Version: "1.2", Size: 10},
			},
		},
		{
			input:  packageTests[2],
			sortBy: packageNameSort,
			expected: []PackageOutput{
				{Name: "a", Version: "1.2", Size: 15},
				{Name: "a", Version: "1.2", Size: 10},
				{Name: "a", Version: "1.4", Size: 20},
			},
		},
	} {
		actual := test.input
		packageBy(test.sortBy).Sort(actual)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("\nExpected: %v\nGot: %v", test.expected, actual)
		}
	}
}

var directoryTests = [][]pkgutil.DirectoryEntry{
	{
		{Name: "a", Size: 10},
		{Name: "b", Size: 12},
		{Name: "c", Size: 20},
	},
	{
		{Name: "a", Size: 10},
		{Name: "b", Size: 12},
		{Name: "c", Size: 12},
	},
}

func TestSortDirectoryEntries(t *testing.T) {
	for _, test := range []struct {
		input    []pkgutil.DirectoryEntry
		sortBy   func(a, b *pkgutil.DirectoryEntry) bool
		expected []pkgutil.DirectoryEntry
	}{
		{
			input:  directoryTests[0],
			sortBy: directorySizeSort,
			expected: []pkgutil.DirectoryEntry{
				{Name: "c", Size: 20},
				{Name: "b", Size: 12},
				{Name: "a", Size: 10},
			},
		},
		{
			input:  directoryTests[0],
			sortBy: directoryNameSort,
			expected: []pkgutil.DirectoryEntry{
				{Name: "a", Size: 10},
				{Name: "b", Size: 12},
				{Name: "c", Size: 20},
			},
		},
		{
			input:  directoryTests[1],
			sortBy: directorySizeSort,
			expected: []pkgutil.DirectoryEntry{
				{Name: "b", Size: 12},
				{Name: "c", Size: 12},
				{Name: "a", Size: 10},
			},
		},
	} {
		actual := test.input
		directoryBy(test.sortBy).Sort(actual)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("\nExpected: %v\nGot: %v", test.expected, actual)
		}
	}
}
