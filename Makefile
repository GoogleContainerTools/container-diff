VERSION_MAJOR ?= 0
VERSION_MINOR ?= 2
VERSION_BUILD ?= 0

GOOS ?= $(shell go env GOOS)
GOARCH = amd64
BUILD_DIR ?= ./out
ORG := github.com/GoogleCloudPlatform
PROJECT := container-diff
REPOPATH ?= $(ORG)/$(PROJECT)

SUPPORTED_PLATFORMS := linux-amd64 darwin-amd64 windows-amd64
BUILD_PACKAGE = $(REPOPATH)

# These build tags are from the containers/image library.
# 
# container_image_ostree_stub allows building the library without requiring the libostree development libraries
# container_image_openpgp forces a Golang-only OpenPGP implementation for signature verification instead of the default cgo/gpgme-based implementation
GO_BUILD_TAGS := "container_image_ostree_stub containers_image_openpgp"
GO_FILES := $(shell go list  -f '{{join .Deps "\n"}}' $(BUILD_PACKAGE) | grep $(ORG) | xargs go list -f '{{ range $$file := .GoFiles }} {{$$.Dir}}/{{$$file}}{{"\n"}}{{end}}')

$(BUILD_DIR)/$(PROJECT): out/$(PROJECT)-$(GOOS)-$(GOARCH)
	cp $(BUILD_DIR)/$(PROJECT)-$(GOOS)-$(GOARCH) $@

$(BUILD_DIR)/$(PROJECT)-%-$(GOARCH): $(GO_FILES) $(BUILD_DIR)
	GOOS=$* GOARCH=$(GOARCH) go build -tags $(GO_BUILD_TAGS) -o $@ $(BUILD_PACKAGE)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: cross
cross: $(foreach platform, $(SUPPORTED_PLATFORMS), out/$(PROJECT)-$(platform))

.PHONY: test
test: $(BUILD_DIR)/$(PROJECT)
	@ ./test.sh

.PHONY: integration
integration: $(BUILD_DIR)/$(PROJECT)
	go test -v -tags integration $(REPOPATH)/tests

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)


