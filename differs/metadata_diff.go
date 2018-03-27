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
	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/GoogleCloudPlatform/container-diff/util"
)

type MetadataAnalyzer struct {
}

type MetadataDiff struct {
	Adds []string
	Dels []string
}

func (a MetadataAnalyzer) Name() string {
	return "MetadataAnalyzer"
}

func (a MetadataAnalyzer) Diff(image1, image2 pkgutil.Image) (util.Result, error) {
	diff, err := getMetadataDiff(image1, image2)
	return &util.MetadataDiffResult{
		Image1:   image1.Source,
		Image2:   image2.Source,
		DiffType: "Metadata",
		Diff:     diff,
	}, err
}

func (a MetadataAnalyzer) Analyze(image pkgutil.Image) (util.Result, error) {
	analysis := getMetadataList(image)
	return &util.ListAnalyzeResult{
		Image:       image.Source,
		AnalyzeType: "Metadata",
		Analysis:    analysis,
	}, nil
}

func getMetadataDiff(image1, image2 pkgutil.Image) (MetadataDiff, error) {
	m1 := getMetadataList(image1)
	m2 := getMetadataList(image2)

	adds := util.GetAdditions(m1, m2)
	dels := util.GetDeletions(m1, m2)
	return MetadataDiff{adds, dels}, nil
}

func getMetadataList(image pkgutil.Image) []string {
	return image.Config.Config.AsList()
}
