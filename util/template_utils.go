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

const FSDiffOutput = `
-----{{.DiffType}}-----

These entries have been added to {{.Image1}}:{{if not .Diff.Adds}} None{{else}}
FILE	SIZE{{range .Diff.Adds}}{{"\n"}}{{.Name}}	{{.Size}}{{end}}{{end}}

These entries have been deleted from {{.Image1}}:{{if not .Diff.Dels}} None{{else}}
FILE	SIZE{{range .Diff.Dels}}{{"\n"}}{{.Name}}	{{.Size}}{{end}}{{end}}

These entries have been changed between {{.Image1}} and {{.Image2}}:{{if not .Diff.Mods}} None{{else}}
FILE	SIZE1	SIZE2{{range .Diff.Mods}}{{"\n"}}{{.Name}}	{{.Size1}}	{{.Size2}}{{end}}
{{end}}
`

const SingleVersionDiffOutput = `
-----{{.DiffType}}-----

Packages found only in {{.Image1}}:{{if not .Diff.Packages1}} None{{else}}
NAME	VERSION	SIZE{{range .Diff.Packages1}}{{"\n"}}{{print "-"}}{{.Name}}	{{.Version}}	{{.Size}}{{end}}{{end}}

Packages found only in {{.Image2}}:{{if not .Diff.Packages2}} None{{else}}
NAME	VERSION	SIZE{{range .Diff.Packages2}}{{"\n"}}{{print "-"}}{{.Name}}	{{.Version}}	{{.Size}}{{end}}{{end}}

Version differences:{{if not .Diff.InfoDiff}} None{{else}}
PACKAGE	IMAGE1 ({{.Image1}})	IMAGE2 ({{.Image2}}){{range .Diff.InfoDiff}}{{"\n"}}{{print "-"}}{{.Package}}	{{.Info1.Version}}, {{.Info1.Size}}	{{.Info2.Version}}, {{.Info2.Size}}{{end}}
{{end}}
`

const MultiVersionDiffOutput = `
-----{{.DiffType}}-----

Packages found only in {{.Image1}}:{{if not .Diff.Packages1}} None{{else}}
NAME	VERSION	SIZE{{range .Diff.Packages1}}{{"\n"}}{{print "-"}}{{.Name}}	{{.Version}}	{{.Size}}{{end}}{{end}}

Packages found only in {{.Image2}}:{{if not .Diff.Packages2}} None{{else}}
NAME	VERSION	SIZE{{range .Diff.Packages2}}{{"\n"}}{{print "-"}}{{.Name}}	{{.Version}}	{{.Size}}{{end}}{{end}}

Version differences:{{if not .Diff.InfoDiff}} None{{else}}
PACKAGE	IMAGE1 ({{.Image1}})	IMAGE2 ({{.Image2}}){{range .Diff.InfoDiff}}{{"\n"}}{{print "-"}}{{.Package}}	{{range .Info1}}{{.Version}}, {{.Size}}{{end}}	{{range .Info2}}{{.Version}}, {{.Size}}{{end}}{{end}}
{{end}}
`

const HistoryDiffOutput = `
-----{{.DiffType}}-----

Docker history lines found only in {{.Image1}}:{{if not .Diff.Adds}} None{{else}}{{block "list" .Diff.Adds}}{{"\n"}}{{range .}}{{print "-" .}}{{"\n"}}{{end}}{{end}}{{end}}

Docker history lines found only in {{.Image2}}:{{if not .Diff.Dels}} None{{else}}{{block "list2" .Diff.Dels}}{{"\n"}}{{range .}}{{print "-" .}}{{"\n"}}{{end}}{{end}}{{end}}
`
const FilenameDiffOutput = `
-----Diff of {{.Filename}}-----
{{.Description}}

{{.Diff}}
`

const ListAnalysisOutput = `
-----{{.AnalyzeType}}-----

Analysis for {{.Image}}:{{if not .Analysis}} None{{else}}{{block "list" .Analysis}}{{"\n"}}{{range .}}{{print "-" .}}{{"\n"}}{{end}}{{end}}{{end}}
`

const FileAnalysisOutput = `
-----{{.AnalyzeType}}-----

Analysis for {{.Image}}:{{if not .Analysis}} None{{else}}
FILE	SIZE{{range .Analysis}}{{"\n"}}{{.Name}}	{{.Size}}{{end}}
{{end}}
`

const MultiVersionPackageOutput = `
-----{{.AnalyzeType}}-----

Packages found in {{.Image}}:{{if not .Analysis}} None{{else}}
NAME	VERSION	SIZE	INSTALLATION{{range .Analysis}}{{"\n"}}{{print "-"}}{{.Name}}	{{.Version}}	{{.Size}}	{{.Path}}{{end}}
{{end}}
`

const SingleVersionPackageOutput = `
-----{{.AnalyzeType}}-----

Packages found in {{.Image}}:{{if not .Analysis}} None{{else}}
NAME	VERSION	SIZE{{range .Analysis}}{{"\n"}}{{print "-"}}{{.Name}}	{{.Version}}	{{.Size}}{{end}}
{{end}}
`
