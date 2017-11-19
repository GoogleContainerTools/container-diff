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

package differs

import (
	"reflect"
	"testing"
)

func TestGetAnalyzers(t *testing.T) {

	tests := []struct {
		name    string
		args    []string
		want    []Analyzer
		wantErr bool
	}{
		{
			name:    "given only one type",
			args:    []string{"history"},
			want:    []Analyzer{HistoryAnalyzer{}},
			wantErr: false,
		},
		{
			name:    "given two type",
			args:    []string{"file", "apt"},
			want:    []Analyzer{FileAnalyzer{}, AptAnalyzer{}},
			wantErr: false,
		},
		{
			name:    "given non-existent type",
			args:    []string{"faketype"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no type given",
			args:    []string{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAnalyzers(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnalyzers() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && tt.wantErr {
				t.Logf("errored out as = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAnalyzers() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
