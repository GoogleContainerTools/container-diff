# container-diff

[![Build
Status](https://travis-ci.org/GoogleCloudPlatform/container-diff.svg?branch=master)](https://travis-ci.org/GoogleCloudPlatform/container-diff)

## What is container-diff?

container-diff is a tool for analyzing and comparing container images. container-diff can examine images along several different criteria, including:
- Docker Image History
- Image file system
- Apt packages
- RPM packages
- pip packages
- npm packages

These analyses can be performed on a single image, or a diff can be performed on two images to compare. The tool can help users better understand what is changing inside their images, and give them a better look at what their images contain.

## Installation

### macOS
```shell
curl -LO https://storage.googleapis.com/container-diff/latest/container-diff-darwin-amd64 && chmod +x container-diff-darwin-amd64 && sudo mv container-diff-darwin-amd64 /usr/local/bin/container-diff
```

### Linux
```shell
curl -LO https://storage.googleapis.com/container-diff/latest/container-diff-linux-amd64 && chmod +x container-diff-linux-amd64 && sudo mv container-diff-linux-amd64 /usr/local/bin/container-diff
```

OR, if you want to avoid using sudo:

```shell
curl -LO https://storage.googleapis.com/container-diff/latest/container-diff-linux-amd64 && chmod +x container-diff-linux-amd64 && mkdir -p $HOME/bin && export PATH=$PATH:$HOME/bin && mv container-diff-linux-amd64 $HOME/bin/container-diff
```

### Windows
Download the [container-diff-windows-amd64.exe](https://storage.googleapis.com/container-diff/latest/container-diff-windows-amd64.exe) file, rename it to `container-diff.exe` and add it to your path


## Quickstart

To use `container-diff analyze` to perform analysis on a single image, you need one Docker image (in the form of an ID, tarball, or URL from a repo). Once you have that image, you can run any of the following analyzers:

```shell
container-diff analyze <img>     [Run default analyzers]
container-diff analyze <img> --type=history  [History]
container-diff analyze <img> --type=file  [File System]
container-diff analyze <img> --type=rpm  [RPM]
container-diff analyze <img> --type=pip  [Pip]
container-diff analyze <img> --type=apt  [Apt]
container-diff analyze <img> --type=node  [Node]
container-diff analyze <img> --type=apt --type=node  [Apt and Node]
# --type=<analyzer1> --type=<analyzer2> --type=<analyzer3>,...
```

By default, with no `--type` flag specified, container-diff will run Apt package analysis.

To use container-diff to perform a diff analysis on two images, you need two Docker images (in the form of an ID, tarball, or URL from a repo). Once you have those images, you can run any of the following differs:
```shell
container-diff diff <img1> <img2>     [Run default differs]
container-diff diff <img1> <img2> --type=history  [History]
container-diff diff <img1> <img2> --type=file  [File System]
container-diff diff <img1> <img2> --type=rpm  [RPM]
container-diff diff <img1> <img2> --type=pip  [Pip]
container-diff diff <img1> <img2> --type=apt  [Apt]
container-diff diff <img1> <img2> --type=node  [Node]
```

You can similarly run many analyzers at once:

```shell
container-diff diff <img1> <img2> --type=history --type=apt --type=node
```

To view the diff of an individual file in two different images, you can use the filename flag in conjuction with the file system diff analyzer.

```shell
container-diff diff <img1> <img2> --type=file --filename=/path/to/file
```

## Image Sources

container-diff supports Docker images located in both a local Docker daemon and a remote registry. To explicitly specify a local image, use the `daemon://` prefix on the image name; similarly, for an explicitly remote image, use the `remote://` prefix.

```shell
container-diff diff daemon://modified_debian:latest remote://gcr.io/google-appengine/debian8:latest
```

Additionally, tarballs can be provided to the tool directly. Make sure your file has a valid tar extension (.tar, .tar.gz, .tgz).

**Note**: container-diff does not support references images by Docker ID directly. If your image only has an ID in your local Docker daemon, you'll need to tag it using `docker tag` before using it with container-diff.

### Authentication

Container-diff supports docker-credential-helpers for authentication when using a registry as an image source.
Make sure you have your credential helper configured before using container-diff, and it should work automatically.

For the Google Container Registry, make sure you have the `docker-credential-gcr` binary configured and on your path, following these [instructions](https://github.com/GoogleCloudPlatform/docker-credential-gcr).


## Other Flags

To get a JSON version of the container-diff output add a `-j` or `--json` flag.

```shell
container-diff diff --type=file --json gcr.io/gcp-runtimes/multi-base gcr.io/gcp-runtimes/multi-modified
```

To order files and packages by size (in descending order) when performing file system or package analyses/diffs, add a `-o` or `--order` flag.

```shell
container-diff analyze remote://gcr.io/gcp-runtimes/multi-modified --type=pip --order
```

To suppress output to stderr, add a `-q` or `--quiet` flag.
```shell
container-diff analyze file1.tar --type=file --quiet
```

## Analysis Result Format

JSON output for analysis results is in the following format:
```json
{
    "Image": "foo",
    "AnalyzeType": "Apt",
    "Analysis": {}
}
```
The possible contents of the `Analysis` field are detailed below.

### History Analysis

The history analyzer outputs a list of strings representing descriptions of how an image layer was created. This is the only analyzer that requires a working Docker daemon to run.

### File System Analysis

The file system analyzer outputs a list of file system contents, including names, paths, and sizes.

### Package Analysis

Package analyzers such as pip, apt, and node inspect the packages installed within the image provided. All package analyses leverage the `PackageOutput` struct, which contains the version and size for a given package instance (and a potential installation path for a specific instance of a package where multiple versions are allowed to be installed), as detailed below:
```go
type PackageOutput struct {
	Name    string
	Path    string
	Version string
	Size    int64
}
```

#### Single Version Package Analysis

Single version package analyzers (apt) have the following output structure: `[]PackageOutput`

Here, the `Path` field is omitted because there is only one instance of each package.

#### Multi Version Package Analysis

Multi version package analyzers (pip, node) have the following output structure: `[]PackageOutput`

Here, the `Path` field is included because there may be more than one instance of each package, and thus the path exists to pinpoint where the package exists in case additional investigation into the package instance is desired.


## Diff Result Format

JSON output for diff results is in the following format:
```json
{
    "Image1": "foo",
    "Image2": "bar",
    "DiffType": "Apt",
    "Diff": {}
}
```
The possible structures of the `Diff` field are detailed below.

### History Diff

The history differ has the following output structure:

```go
type HistDiff struct {
	Adds []string
	Dels []string
}
```

### File System Diff

The file system differ has the following output structure:

```go
type DirDiff struct {
	Adds  []string
	Dels  []string
	Mods  []string
}
```

### Package Diffs

Package differs such as pip, apt, and node inspect the packages contained within the images provided. All packages differs currently leverage the PackageInfo struct which contains the version and size for a given package instance, as detailed below:
```go
type PackageInfo struct {
	Version string
	Size	string
}
```

#### Single Version Package Diffs

Single version differs (apt) have the following JSON output structure:

```go
type PackageDiff struct {
	Packages1 []PackageOutput
	Packages2 []PackageOutput
	InfoDiff  []Info
}
```

Packages1 and Packages2 detail which packages exist uniquely in Image1 and Image2, respectively, with package name, version and size info. InfoDiff contains a list of Info structs, each of which contains the package name (which occurred in both images but had a difference in size or version), and the PackageInfo struct for each package instance.

#### Multi Version Package Diffs

The multi version differs (pip, node) support processing images which may have multiple versions of the same package. Below is the json output structure:

```go
type MultiVersionPackageDiff struct {
	Packages1 []PackageOutput
	Packages2 []PackageOutput
	InfoDiff  []MultiVersionInfo
}
```

Packages1 and Packages2 detail which packages exist uniquely in Image1 and Image2, respectively, with package name, installation path, version and size info.  InfoDiff here is exanded to allow for multiple versions to be associated with a single package. In this case, a package of the same name is considered to differ between two images when there exist one or more instances of it installed in one image but not the other (i.e. have a unique version and/or size).

```go
type MultiVersionInfo struct {
	Package string
	Info1	[]PackageInfo
	Info2	[]PackageInfo
}
```

## User Customized Output
Users can customize the format of the output of diffs with the`--format` flag. The flag takes a Go template string, which specifies the format the diff should be output in. This template string uses the structs described above, depending on the differ used, to format output.  The default template strings container-diff uses can be found [here](https://github.com/GoogleCloudPlatform/container-diff/blob/master/util/template_utils.go).

An example using the pip package analyzer is shown below, in which only package names are printed (some are repeated because of version differences).

```shell
$ container-diff analyze gcr.io/google-appengine/python:latest --type=pip --format='
-----{{.AnalyzeType}}-----
Packages found in {{.Image}}:{{if not .Analysis}} None{{else}}
{{range .Analysis}}{{"\n"}}{{.Name}}{{end}}
{{end}}
'
Retrieving image gcr.io/google-appengine/python:latest from source Cloud Registry
Retrieving analyses

-----Pip-----
Packages found in gcr.io/google-appengine/python:latest:

chardet
colorama
html5lib
mercurial
pip
pip
pip
requests
setuptools
setuptools
setuptools
six
urllib3
virtualenv
wheel
wheel
```

## Known issues

To run container-diff using image IDs, docker must be installed.
Tarballs provided directly to the tool must be in the Docker format (i.e. have a manifest.json file for layer ordering)


## Example Run

```shell
$ container-diff diff gcr.io/google-appengine/python:2017-07-21-123058 gcr.io/google-appengine/python:2017-06-29-190410 --type=apt --type=node --type=pip

-----AptDiffer-----

Packages found only in gcr.io/google-appengine/python:2017-07-21-123058: None

Packages found only in gcr.io/google-appengine/python:2017-06-29-190410: None

Version differences:
PACKAGE             IMAGE1 (gcr.io/google-appengine/python:2017-07-21-123058)        IMAGE2 (gcr.io/google-appengine/python:2017-06-29-190410)
-libgcrypt20        1.6.3-2 deb8u4, 998K                                             1.6.3-2 deb8u3, 1002K

-----NodeDiffer-----

Packages found only in gcr.io/google-appengine/python:2017-07-21-123058: None

Packages found only in gcr.io/google-appengine/python:2017-06-29-190410: None

Version differences: None

-----PipDiffer-----

Packages found only in gcr.io/google-appengine/python:2017-07-21-123058: None

Packages found only in gcr.io/google-appengine/python:2017-06-29-190410: None

Version differences: None

```

```shell
$ container-diff diff file1.tar file2.tar --type=file --filename=go/src/app/file.txt
Starting diff on images file1.tar and file2.tar, using differs: [file]
Retrieving image file2.tar from source Tar Archive
Retrieving image file1.tar from source Tar Archive
Computing diffs

-----File-----

These entries have been added to file1.tar: None

These entries have been deleted from file1.tar: None

These entries have been changed between file1.tar and file2.tar:
FILE                        SIZE1        SIZE2
/go/src/app/file.txt        30B          30B

Computing filename diffs

-----Diff of go/src/app/file.txt-----


--- file1.tar
+++ file2.tar
@@ -1 +1 @@
-This is file 1
This is a file
+This is file 2
This is a file
```

## Example Run with JSON post-processing
The following example demonstrates how one might selectively display the output of their diff, such that version differences are ignored and only package absence/presence is displayed and the packages present in only one image are sorted by size in descending order. A small piece of the JSON being post-processed can be seen below:
```json
[
  {
    "DiffType": "AptDiffer",
    "Diff": {
      "Image1": "gcr.io/gcp-runtimes/multi-base",
      "Packages1": {},
      "Image2": "gcr.io/gcp-runtimes/multi-modified",
      "Packages2": {
        "dh-python": {
          "Version": "1.20141111-2",
          "Size": "277"
        },
        "libmpdec2": {
          "Version": "2.4.1-1",
          "Size": "275"
        }
      }
    }
  }
]
```
The post-processing script used for this example is below:

```python
import sys, json

def main():
  data = json.loads(sys.stdin.read())
  img1packages = []
  img2packages = []
  for differ in data:
    diff = differ['Diff']

    if len(diff['Packages1']) > 0:
      for package in diff['Packages1']:
        Size = package['Size']
        img1packages.append((str(package), int(str(Size))))

    if len(diff['Packages2']) > 0:
      for package in diff['Packages2']:
        Size = package['Size']
        img2packages.append((str(package), int(str(Size))))

    img1packages = reversed(sorted(img1packages, key=lambda x: x[1]))
    img2packages = reversed(sorted(img2packages, key=lambda x: x[1]))


    print "Only in image1\n"
    for pkg in img1packages:
      print pkg
    print "Only in image2\n"
    for pkg in img2packages:
      print pkg
    print

if __name__ == "__main__":
  main()
```

Given the above python script to postprocess json output, you can produce the following behavior:
```shell
container-diff gcr.io/gcp-runtimes/multi-base gcr.io/gcp-runtimes/multi-modified -a -j | python pyscript.py

Only in image1

Only in image2

('libpython3.4-stdlib', 9484)
('python3.4-minimal', 4506)
('libpython3.4-minimal', 3310)
('python3.4', 336)
('dh-python', 277)
('libmpdec2', 275)
('python3-minimal', 96)
('python3', 36)
('libpython3-stdlib', 28)

```
## Make your own differ

Feel free to develop your own analyzer leveraging the utils currently available. PRs are welcome!

### Custom Analyzer Quickstart

In order to quickly make your own analyzer, follow these steps:

1. Determine if you can use existing analyzing or diffing tools.  If you can make use of existing tools, you then need to construct the structs to feed into the tools by getting all of the packages for each image or the analogous quality to be analyzed.  To determine if you can leverage existing tools, think through these questions:
- Are you trying to analyze packages?
    - Yes: Does the relevant package manager support different versions of the same package on one image?
        - Yes: Implement `getPackages` to collect all versions of all packages within an image in a `map[string]map[string]util.PackageInfo`. Use [`GetMultiVersionMapDiff`](https://github.com/GoogleCloudPlatform/container-diff/blob/0031c88993c9ac019e2d404815ef50c652d8d010/util/package_diff_utils.go#L119-L126) to diff map objects.  See [`differs/node_diff.go`](https://github.com/GoogleCloudPlatform/container-diff/blob/0031c88993c9ac019e2d404815ef50c652d8d010/differs/node_diff.go#L49-L93) or [`differs/pip_diff.go`](https://github.com/GoogleCloudPlatform/container-diff/blob/0031c88993c9ac019e2d404815ef50c652d8d010/differs/pip_diff.go#L48-L111) for examples.
        -  No: Implement `getPackages` to collect all versions of all packages within an image in a `map[string]util.PackageInfo`. Use [`GetMapDiff`](https://github.com/GoogleCloudPlatform/container-diff/blob/31cec2304b54ae6ae444ccde4382b113d8e06097/util/package_diff_utils.go#L110-L117) to diff map objects.  See [`differs/apt_diff.go`](https://github.com/GoogleCloudPlatform/container-diff/blob/master/differs/apt_diff.go#L29).
    - No: Look to [History](https://github.com/GoogleCloudPlatform/container-diff/blob/0031c88993c9ac019e2d404815ef50c652d8d010/differs/history_diff.go) and [File System](https://github.com/GoogleCloudPlatform/container-diff/blob/0031c88993c9ac019e2d404815ef50c652d8d010/differs/file_diff.go) differs as models for diffing.

2. Write your analyzer driver in the `differs` directory, such that you have a struct for your analyzer type and two methods for that analyzer: `Analyze` for single image analysis and `Diff` for comparison between two images:

```go
type YourAnalyzer struct {}

func (a YourAnalyzer) Analyze(image util.Image) (util.Result, error) {...}
func (a YourAnalyzer) Diff(image1, image2 util.Image) (util.Result, error) {...}
```
The image arguments passed to your analyzer contain the path to the unpacked tar representation of the image, as well as certain configuration information (e.g. environment variables upon image creation and image history).

If using existing package tools, you should create the appropriate structs (e.g. `SingleVersionPackageAnalyzeResult` or `SingleVersionPackageDiffResult`) to analyze or diff.  Otherwise, create your own structs which should yield information to fill an AnalyzeResult or DiffResult as the return type for Analyze() and Diff(), respectively, and should implement the `Result` interface, as in the next step.

3. Create a struct following the [`Result`](https://github.com/GoogleCloudPlatform/container-diff/blob/0031c88993c9ac019e2d404815ef50c652d8d010/util/analyze_output_utils.go#L27-L30) interface by implementing the following two methods.

```go
type Result interface {
	OutputStruct() interface{}
	OutputText(resultType string, format string) error
}
```

This is where you define how your analyzer should output for a human readable format (`OutputText`) and as a struct which can then be written to a `.json` file.  See [`util/diff_output_utils.go`](https://github.com/GoogleCloudPlatform/container-diff/blob/0031c88993c9ac019e2d404815ef50c652d8d010/util/diff_output_utils.go) and [`util/analyze_output_utils.go`](https://github.com/GoogleCloudPlatform/container-diff/blob/0031c88993c9ac019e2d404815ef50c652d8d010/util/analyze_output_utils.go).

4. Add your analyzer to the `Analyzers` map in [`differs/differs.go`](https://github.com/GoogleCloudPlatform/container-diff/blob/0031c88993c9ac019e2d404815ef50c652d8d010/differs/differs.go#L44-L50) with the corresponding Analyzer struct as the value.
