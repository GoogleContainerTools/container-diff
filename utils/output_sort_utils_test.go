package utils

import (
	"reflect"
	"testing"
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

var directoryTests = [][]DirectoryEntry{
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
		input    []DirectoryEntry
		sortBy   func(a, b *DirectoryEntry) bool
		expected []DirectoryEntry
	}{
		{
			input:  directoryTests[0],
			sortBy: directorySizeSort,
			expected: []DirectoryEntry{
				{Name: "c", Size: 20},
				{Name: "b", Size: 12},
				{Name: "a", Size: 10},
			},
		},
		{
			input:  directoryTests[0],
			sortBy: directoryNameSort,
			expected: []DirectoryEntry{
				{Name: "a", Size: 10},
				{Name: "b", Size: 12},
				{Name: "c", Size: 20},
			},
		},
		{
			input:  directoryTests[1],
			sortBy: directorySizeSort,
			expected: []DirectoryEntry{
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
