package cmd

import (
	"testing"
)

var analyzeArgNumTests = []testpair{
	{[]string{}, false},
	{[]string{"one"}, true},
	{[]string{"one", "two"}, false},
}

func TestAnalyzeArgNum(t *testing.T) {
	for _, test := range analyzeArgNumTests {
		err := checkAnalyzeArgNum(test.input)
		if (err == nil) != test.expected_output {
			if test.expected_output {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				t.Errorf("Expected error but got none")
			}
		}
	}
}
