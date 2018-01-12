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
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/sirupsen/logrus"
)

var templates = map[string]string{
	"SingleVersionPackageDiff":    SingleVersionDiffOutput,
	"MultiVersionPackageDiff":     MultiVersionDiffOutput,
	"HistDiff":                    HistoryDiffOutput,
	"DirDiff":                     FSDiffOutput,
	"FilenameDiff":                FilenameDiffOutput,
	"ListAnalyze":                 ListAnalysisOutput,
	"FileAnalyze":                 FileAnalysisOutput,
	"MultiVersionPackageAnalyze":  MultiVersionPackageOutput,
	"SingleVersionPackageAnalyze": SingleVersionPackageOutput,
}

func JSONify(diff interface{}) error {
	diffBytes, err := json.MarshalIndent(diff, "", "  ")
	if err != nil {
		return err
	}
	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	f.Write(diffBytes)
	return nil
}

func getTemplate(templateType string) (string, error) {
	if template, ok := templates[templateType]; ok {
		return template, nil
	}
	return "", errors.New("No available template")
}

func TemplateOutput(diff interface{}, templateType string) error {
	outputTmpl, err := getTemplate(templateType)
	if err != nil {
		logrus.Error(err)

	}
	funcs := template.FuncMap{"join": strings.Join}
	tmpl, err := template.New("tmpl").Funcs(funcs).Parse(outputTmpl)
	if err != nil {
		logrus.Error(err)
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 8, ' ', 0)
	err = tmpl.Execute(w, diff)
	if err != nil {
		logrus.Error(err)
		return err
	}
	w.Flush()
	return nil
}

func TemplateOutputFromFormat(diff interface{}, templateType string, format string) error {
	if format == "" {
		return TemplateOutput(diff, templateType)
	}
	funcs := template.FuncMap{"join": strings.Join}
	tmpl, err := template.New("tmpl").Funcs(funcs).Parse(format)
	if err != nil {
		logrus.Warningf("User specified format resulted in error, printing default output.")
		logrus.Error(err)
		return TemplateOutput(diff, templateType)
	}
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 8, ' ', 0)
	err = tmpl.Execute(w, diff)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}
