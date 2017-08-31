package cmd

import (
	"testing"
)

type testpair struct {
	input           []string
	expected_output bool
}

var argTypeTests = []testpair{
	{[]string{"badID", "badID"}, false},
	{[]string{"123456789012", "badID"}, false},
	{[]string{"123456789012", "123456789012"}, true},
	{[]string{"?!badDiffer71", "123456789012"}, false},
	{[]string{"123456789012", "gcr.io/repo/image"}, true},
}

func TestArgType(t *testing.T) {
	for _, test := range argTypeTests {
		err := checkArgType(test.input)
		if (err == nil) != test.expected_output {
			if test.expected_output {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				t.Errorf("Expected error but got none")
			}
		}
	}
}
