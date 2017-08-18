# container-diff

[![Build
Status](https://travis-ci.org/GoogleCloudPlatform/container-diff.svg?branch=master)](https://travis-ci.org/GoogleCloudPlatform/container-diff)

## What is container-diff?

container-diff is an image differ command line tool.  container-diff can diff two images along several different criteria, currently including:
- Docker Image History
- Image file system
- apt-get installed packages
- pip installed packages
- npm installed packages

This tool can help you as a developer better understand what is changing within your images and better understand what your images contain.

## Installation

### macOS
```shell
curl -LO container-diff https://storage.googleapis.com/container-diff/v0.1.0/container-diff-darwin-amd64
```

### Linux
```shell
curl -LO https://storage.googleapis.com/container-diff/v0.1.0/container-diff-linux-amd64 && chmod +x container-diff-linux-amd64 && sudo mv container-diff-linux-amd64 /usr/local/bin/
```

### Windows
Download the [container-diff-windows-amd64.exe](https://storage.googleapis.com/container-diff/v0.1.0/container-diff-windows-amd64.exe) file, rename it to `container-diff.exe` and add it to your path


## Quickstart

To use container-diff you need two Docker images (in the form of an ID, tarball, or URL from a repo).  Once you have those images you can run any of the following differs:

```
container-diff <img1> <img2>     [Run all differs]
container-diff <img1> <img2> -d  [History]
container-diff <img1> <img2> -f  [File System]
container-diff <img1> <img2> -p  [Pip]
container-diff <img1> <img2> -a  [Apt]
container-diff <img1> <img2> -n  [Node]
```

You can similarly run many differs at once:

```
container-diff <img1> <img2> -d -a -n [History, Apt, and Node]
```
All of the differ flags with their long versions can be seen below:

| Differ                    | Short flag | Long Flag  |
| ------------------------- |:----------:| ----------:|
| File System diff          | -f         | --file     |
| History                   | -d 	 | --history  |
| npm installed packages    | -n 	 | --node     |
| pip installed packages    | -p 	 | --pip      |
| apt-get installed packages| -a 	 | --apt      |




## Other Flags

To get a JSON version of the container-diff output add a `-j` or `--json` flag.

```container-diff <img1> <img2> -j```

To use the docker client instead of shelling out to your local docker daemon, add a `-e` or `--eng` flag.

```container-diff <img1> <img2> -e```


## Output Format

### History Diff

The history differ has the following json output structure:

```
type HistDiff struct {
	Image1 string
	Image2 string
	Adds   []string
	Dels   []string
}
```

### File System Diff

The files system differ has the following json output structure: 

```
type DirDiff struct {
	Image1 string
	Image2 string
	Adds   []string
	Dels   []string
	Mods   []string
}
```

### Package Diffs

Package differs such as pip, apt, and node inspect the packages contained within the images provided.  All packages differs currently leverage the PackageInfo struct which contains the version and size for a given package instance.

```
type PackageInfo struct {
	Version string
	Size    string
}
```

#### Single Version Diffs

Single version differs (apt) have the following json output structure:

```
type PackageDiff struct {
	Image1    string
	Packages1 map[string]PackageInfo
	Image2    string
	Packages2 map[string]PackageInfo
	InfoDiff  []Info
}
```

Image1 and Image2 are the image names.  Packages1 and Packages2 map package names to PackageInfo structs which contain the version and size of the package.  InfoDiff contains a list of Info structs, each of which contains the package name (which occurred in both images but had a difference in size or version), and the PackageInfo struct for each package instance. 

#### Multi Version Diffs

The multi version differs (pip, node) support processing images which may have multiple versions of the same package.  Below is the json output structure:

```
type MultiVersionPackageDiff struct {
	Image1    string
	Packages1 map[string]map[string]PackageInfo
	Image2    string
	Packages2 map[string]map[string]PackageInfo
	InfoDiff  []MultiVersionInfo
}
```

Image1 and Image2 are the image names.  Packages1 and Packages2 map package name to path where the package was found to PackageInfo struct (version and size of that package instance).  InfoDiff here is exanded to allow for multiple versions to be associated with a single package.

```
type MultiVersionInfo struct {
	Package string
	Info1   []PackageInfo
	Info2   []PackageInfo
}
```

## Known issues

To run container-diff on image IDs, docker must be installed.

If encountering this error ```open /etc/docker/certs.d/gcr.io: permission
denied```, run ```sudo rm -rf /etc/docker```.

## Example Run

```
$ container-diff gcr.io/google-appengine/python:2017-07-21-123058 gcr.io/google-appengine/python:2017-06-29-190410 -a -n -p

-----AptDiffer-----

Packages found only in gcr.io/google-appengine/python:2017-07-21-123058: None

Packages found only in gcr.io/google-appengine/python:2017-06-29-190410: None

Version differences:
PACKAGE             IMAGE1 (gcr.io/google-appengine/python:2017-07-21-123058)        IMAGE2 (gcr.io/google-appengine/python:2017-06-29-190410)
-libgcrypt20        1.6.3-2 deb8u4, 998B                                             1.6.3-2 deb8u3, 1002B

-----NodeDiffer-----

Packages found only in gcr.io/google-appengine/python:2017-07-21-123058: None

Packages found only in gcr.io/google-appengine/python:2017-06-29-190410: None

Version differences: None

-----PipDiffer-----

Packages found only in gcr.io/google-appengine/python:2017-07-21-123058: None

Packages found only in gcr.io/google-appengine/python:2017-06-29-190410: None

Version differences: None

```
## Example Run with json post-processing
The following example demonstrates how one might selectively display the output of their diff, such that version differences are ignored and only package absence/presence is displayed and the packages present in only one image are sorted by size in descending order.  A small piece of the json being post-processed can be seen below:
```
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
        },
        ...
```
The post-processing script used for this example is below:

```import sys, json

def main():
  data = json.loads(sys.stdin.read())
  img1packages = []
  img2packages = []
  for differ in data:
    diff = differ['Diff']

    if len(diff['Packages1']) > 0:
      for package in diff['Packages1']:
        Size = diff['Packages1'][package]['Size']
        img1packages.append((str(package), int(str(Size))))

    if len(diff['Packages2']) > 0:
      for package in diff['Packages2']:
        Size = diff['Packages2'][package]['Size']
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
```
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

Feel free to develop your own differ leveraging the utils currently available.  PRs are welcome.

### Custom Differ Quickstart

In order to quickly make your own differ, follow these steps:

1. Add your diff identifier to the flags in [root.go](https://github.com/GoogleCloudPlatform/container-diff/blob/ReadMe/cmd/root.go)
2. Determine if you can use existing differ tools.  If you can make use of existing tools, you then need to construct the structs to feed into the diff tools by getting all of the packages for each image or the analogous quality to be diffed.  To determine if you can leverage existing tools, think through these questions:
- Are you trying to diff packages?
    - Yes: Does the relevant package manager support different versions of the same package on one image?
        - Yes: Use `GetMultiVerisonMapDiff` to diff `map[string]map[string]utils.PackageInfo` objects.  See [nodeDiff.go](https://github.com/GoogleCloudPlatform/container-diff/blob/master/differs/nodeDiff.go#L33)  or [pipDiff.go](https://github.com/GoogleCloudPlatform/container-diff/blob/master/differs/pipDiff.go#L23) for examples.
        -  No: Use `GetMapDiff` to diff `map[string]utils.PackageInfo` objects.  See [aptDiff.go](https://github.com/GoogleCloudPlatform/container-diff/blob/master/differs/aptDiff.go#L29). 
    - No: Look to [History](https://github.com/GoogleCloudPlatform/container-diff/blob/ReadMe/differs/historyDiff.go) and [File System](https://github.com/GoogleCloudPlatform/container-diff/blob/ReadMe/differs/fileDiff.go) differs as models for diffing.

3. Write your Diff driver in the `differs` directory, such that you have a struct for your differ type and a method for that differ called Diff:

```
type YourDiffer struct {}

func (d YourDiffer) Diff(image1, image2 utils.Image) (DiffResult, error) {...}
```
The arguments passed to your differ contain the path to the unpacked tar representation of the image.  That path can be accessed as such: `image1.FSPath`.  

If using existing package differ tools, you should create the appropriate structs to diff (determined in step 2 - either `map[string]map[string]utils.PackageInfo` or `map[string]utils.PackageInfo`) and then call the appropriate get diff function (also determined in step2 - either `GetMultiVerisonMapDiff` or `GetMapDiff`).

Otherwise, create your own differ which should yield information to fill a DiffResult in the next step.

4. Create a DiffResult for your differ.  
```
type DiffResult interface {
	GetStruct() DiffResult
	OutputText(diffType string) error
}
```

This is where you define how your differ should output for a human readable format (`OutputText`) and as a struct which can then be written to a `.json` file.  See [output_utils.go](https://github.com/GoogleCloudPlatform/container-diff/blob/master/utils/output_utils.go).

5. Add your differ to the diffs map in [differs.go](https://github.com/GoogleCloudPlatform/container-diff/blob/master/differs/differs.go#L22) with the corresponding Differ struct as the value.





