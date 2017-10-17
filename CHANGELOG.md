# container-diff Release Notes

# Version 0.5.0 - 10/9/17
* Apt diffing now done by default
* Add support for building single platform with Bazel
* Move util methods to new `pkgutil` package for vendoring
* Add support for specifying local vs remote images with `daemon://` and `remote://` prefixes
* Remove Docker dependency for local images


# Version 0.4.1
* Fixed error with running container-diff with no analyzer specified
* Fixed error where container-diff version was outputting an incorrect value

# Version 0.4.0
* Added single image analysis #20 
* Added file/package output sorting by size #36
* Changed CLI to use "--types" flag #68 
* Various cleaning and refactoring
