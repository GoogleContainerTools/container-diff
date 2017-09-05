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

package utils

import (
	"github.com/pmezard/go-difflib/difflib"
)

// Modification of difflib's unified differ
func GetAdditions(a, b []string) []string {
	matcher := difflib.NewMatcher(a, b)
	differences := matcher.GetGroupedOpCodes(0)

	adds := []string{}
	for _, group := range differences {
		for _, opCode := range group {
			j1, j2 := opCode.J1, opCode.J2
			if opCode.Tag == 'r' || opCode.Tag == 'i' {
				for _, line := range b[j1:j2] {
					adds = append(adds, line)
				}
			}
		}
	}
	return adds
}

func GetDeletions(a, b []string) []string {
	matcher := difflib.NewMatcher(a, b)
	differences := matcher.GetGroupedOpCodes(0)

	dels := []string{}
	for _, group := range differences {
		for _, opCode := range group {
			i1, i2 := opCode.I1, opCode.I2
			if opCode.Tag == 'r' || opCode.Tag == 'd' {
				for _, line := range a[i1:i2] {
					dels = append(dels, line)
				}
			}
		}
	}
	return dels
}

func GetMatches(a, b []string) []string {
	matcher := difflib.NewMatcher(a, b)
	matchindexes := matcher.GetMatchingBlocks()

	matches := []string{}
	for i, match := range matchindexes {
		if i != len(matchindexes)-1 {
			start := match.A
			end := match.A + match.Size
			for _, line := range a[start:end] {
				matches = append(matches, line)
			}
		}
	}
	return matches
}
