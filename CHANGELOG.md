# v0.17.0 Release - 06/07/2021
**Linux**
`curl -LO https://storage.googleapis.com/container-diff/v0.17.0/container-diff-linux-amd64 && mv container-diff-linux-amd64 container-diff && chmod +x container-diff && sudo mv container-diff /usr/local/bin/`

**macOS**
`curl -LO https://storage.googleapis.com/container-diff/v0.17.0/container-diff-darwin-amd64 && mv container-diff-darwin-amd64 container-diff && chmod +x container-diff && sudo mv container-diff /usr/local/bin/`

**Windows**
https://storage.googleapis.com/container-diff/v0.17.0/container-diff-windows-amd64.exe


**Note from the maintainers**: container-diff has been moved to maintenance mode, but this does NOT mean the project is shutting down! Unfortunately, our team at Google does not have the bandwidth to actively maintain this project, and we want to be sure that the OSS community knows that we're not ignoring our users and contributors. We'll continue to provide critical fixes, and review and ship new contributions as they're submitted to the project.

Fixes:
* Resolve symlink issue. [#355](https://github.com/GoogleContainerTools/container-diff/pull/355)
* Remove global variable [#349](https://github.com/GoogleContainerTools/container-diff/pull/349)
* Code improvement in function readErrorsFromChannel() [#347](https://github.com/GoogleContainerTools/container-diff/pull/347)

Huge thank you to all of our dedicated contributors:
- Leonor Resende
- Nick Kubala
- nicolasdilley
- zhongjie


# v0.16.0 Release - 12/21/2020
**Linux**
`curl -LO https://storage.googleapis.com/container-diff/v0.16.0/container-diff-linux-amd64 && mv container-diff-linux-amd64 container-diff && chmod +x container-diff && sudo mv container-diff /usr/local/bin/`

**macOS**
`curl -LO https://storage.googleapis.com/container-diff/v0.16.0/container-diff-darwin-amd64 && mv container-diff-darwin-amd64 container-diff && chmod +x container-diff && sudo mv container-diff /usr/local/bin/`

**Windows**
https://storage.googleapis.com/container-diff/v0.16.0/container-diff-windows-amd64.exe


**Note from the maintainers**: container-diff has been moved to maintenance mode, but this does NOT mean the project is shutting down! Unfortunately, our team at Google does not have the bandwidth to actively maintain this project, and we want to be sure that the OSS community knows that we're not ignoring our users and contributors. We'll continue to provide critical fixes, and review and ship new contributions as they're submitted to the project.

Highlights:
* container-diff now supports packages installed with [Emerge](https://wiki.gentoo.org/wiki/Portage)!

New Features:
* feat: support emerge packages analyzer [#337](https://github.com/GoogleContainerTools/container-diff/pull/337)
* Add two options to handle self-signed certificates registries [#327](https://github.com/GoogleContainerTools/container-diff/pull/327)

Fixes:
* version: Move vX.Y.Z to version.go so it works with `go get`, add git info [#304](https://github.com/GoogleContainerTools/container-diff/pull/304)
* Fix RPM differ to to include release of version [#315](https://github.com/GoogleContainerTools/container-diff/pull/315)
* --help: List available analyzers, improve Usage line [#303](https://github.com/GoogleContainerTools/container-diff/pull/303)
* Fix concurrent map write for hardlink [#324](https://github.com/GoogleContainerTools/container-diff/pull/324)
* Remove unnecessary flag parsing [#330](https://github.com/GoogleContainerTools/container-diff/pull/330)

Updates:
* Upgrade to go 1.14 and go.mod [#329](https://github.com/GoogleContainerTools/container-diff/pull/329)
* Update codeowners [#341](https://github.com/GoogleContainerTools/container-diff/pull/341)

Docs Updates:
* README: mention archlinux-specific instructions [#307](https://github.com/GoogleContainerTools/container-diff/pull/307)
* Document official support level from Google [#342](https://github.com/GoogleContainerTools/container-diff/pull/342)

Huge thanks goes out to all of our contributors for this release:

- Ben Einaudi
- Beni Cherniavsky-Paskin
- Callum Reardon
- Don McCasland
- Lexus Lee
- Luis Plazas
- Nick Kubala
- Santiago Torres


# Version 0.15.0 - 02/19/19
* Update deps [#298](https://github.com/GoogleContainerTools/container-diff/pull/298)
* Fix result switch while viewing with type history [#296](https://github.com/GoogleContainerTools/container-diff/pull/296)
* Use PKG-INFO and METADATA to infer package names in pip analysis [#292](https://github.com/GoogleContainerTools/container-diff/pull/292)
* Remove characters from cache path that are invalid on Windows [#285](https://github.com/GoogleContainerTools/container-diff/pull/285)
* Use top_level.txt when analyzing pip modules [#291](https://github.com/GoogleContainerTools/container-diff/pull/291)
* Strip colons from file path before creating cache dir [#290](https://github.com/GoogleContainerTools/container-diff/pull/290)
* Adding Github Action to run Container Diff [#286](https://github.com/GoogleContainerTools/container-diff/pull/286)

# Version 0.14.0 - 12/18/18
* Enhancement - save to file [#279](https://github.com/GoogleContainerTools/container-diff/pull/279)
* Fixed concurrent map write in image diffing [#278](https://github.com/GoogleContainerTools/container-diff/pull/278)
* Adding custom cache envar and command line argument [#274](https://github.com/GoogleContainerTools/container-diff/pull/274)
* Split lines prior to diffing [#272](https://github.com/GoogleContainerTools/container-diff/pull/272)
* Move all image processing logic into utils, and expose publically [#270](https://github.com/GoogleContainerTools/container-diff/pull/270)
* Enhancement - save to file [#279](https://github.com/GoogleContainerTools/container-diff/pull/279)
* Fixed concurrent map write in image diffing [#278](https://github.com/GoogleContainerTools/container-diff/pull/278)
* Adding custom cache envar and command line argument [#274](https://github.com/GoogleContainerTools/container-diff/pull/274)
* Split lines prior to diffing [#272](https://github.com/GoogleContainerTools/container-diff/pull/272)
* Move all image processing logic into utils, and expose publically [#270](https://github.com/GoogleContainerTools/container-diff/pull/270)

# Version 0.13.0 - 10/30/18
* Update go-containerregistry to pick up docker API client negotiation [#267](https://github.com/GoogleContainerTools/container-diff/pull/267)
* Fix unintended variable shadowing [#263](https://github.com/GoogleContainerTools/container-diff/pull/263)
* Change the default analysis type from apt to size [#266](https://github.com/GoogleContainerTools/container-diff/pull/266)

# Version 0.12.0 - 10/1/18
* Add script to list all pull requests for each release [#258](https://github.com/GoogleContainerTools/container-diff/pull/258)
* Fix deps [#260](https://github.com/GoogleContainerTools/container-diff/pull/260)
* Backfill changelog [#257](https://github.com/GoogleContainerTools/container-diff/pull/257)
* Add maintainers file and new issue template [#259](https://github.com/GoogleContainerTools/container-diff/pull/259)
* Add size analyzer [#256](https://github.com/GoogleContainerTools/container-diff/pull/256)
* Fix destination path for clone in contrib guidance. [#255](https://github.com/GoogleContainerTools/container-diff/pull/255)
* Add rpmlayer differ [#252](https://github.com/GoogleContainerTools/container-diff/pull/252)
* Handle error gracefully when we can't retrieve an image [#251](https://github.com/GoogleContainerTools/container-diff/pull/251)
* Layered analysis for single version packages [#248](https://github.com/GoogleContainerTools/container-diff/pull/248)
* Reuse cached filesystems for layers [#247](https://github.com/GoogleContainerTools/container-diff/pull/247)

# Version 0.11.0 - 6/27/18
* Don't overwrite loaded tarball image
* Use local RPM binary (when possible) in RPM differ

# Version 0.10.0 - 6/13/18
* Switch to github.com/google/go-containerregistry
* Fix entrypoint in RPM differ
* Various metadata diffing fixes
* Remove Bazel

# Version 0.9.0 - 4/10/18
* Add metadata diffing
* Sanitize filepaths before joining to prevent filepath traversal
* Fix appending of latest tag to tar files
* Correctly clean up image filesystems
* Set/unset write bit when unpacking directories out of permission scope
* Add all docker config fields to image config
* Various bug/panic fixes

# Version 0.8.0 - 3/19/28
* Allow updating env vars on MutableSource image
* Save temp layers in cache directory instead of /tmp
* Allow accessing and modifying MutableSource config
* Fixed appending latest tag to images with no tag provided
* Created default ImageSource if none is provided to prepper
* Fixed issue where remote:// prefix was not being stripped correctly

# Version 0.7.0 - 2/22/18
* Download remote:// images in RPMAnalyzer
* Add support for custom formatting strings
* Refactors to the cache and image unpacking code
* Add Label to ConfigObject
* Add MutableSource for basic image modifications

# Version 0.6.2 - 1/10/18
* Fix issue with user.Current not working in some environments

# Version 0.6.1 - 1/4/18
* Fix incorrect version in binary

# Version 0.6.0 - 12/27/17
* Add support for diffing RPM packages
* Fix a few unpack errors for images with whiteout layers
* Switch dependency management from godep to dep

# Version 0.5.2 - 11/25/17
* Various docs fixes
* Fix Makefile to preserve all build artifacts

# Version 0.5.1 - 11/20/17
* Change types flag from comma separated --types list to repeated --typeflag
* Added --filename flag to show diffs of individual files
* Added layer caching

# Version 0.5.0 - 10/9/17
* Apt diffing now done by default
* Add support for building single platform with Bazel
* Move util methods to new `pkgutil` package for vendoring
* Add support for specifying local vs remote images with `daemon://` and `remote://` prefixes
* Remove Docker dependency for local images


# Version 0.4.1 - 9/12/17
* Fixed error with running container-diff with no analyzer specified
* Fixed error where container-diff version was outputting an incorrect value

# Version 0.4.0 - 9/12/17
* Added single image analysis #20 
* Added file/package output sorting by size #36
* Changed CLI to use "--types" flag #68 
* Various cleaning and refactoring
