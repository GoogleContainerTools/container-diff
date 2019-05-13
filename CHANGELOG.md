# container-diff Release Notes

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
