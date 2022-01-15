build_dir := build

name_cmd := jiratool

pkg_base := bitbucket.org/christian_m/jiratool
pkg_cmd := $(pkg_base)/cmd/$(name_cmd)

version := $(if $(shell git describe --tags --abbrev=0),$(shell git describe --tags --abbrev=0),build_$(shell git rev-parse --short HEAD))

build = GOOS=$(1) GOARCH=$(2) go build -o $(build_dir)/$(3)$(4) $(5)
tar = cd $(build_dir) && tar -cvzf $(2)_$(1)_$(version).tar.gz $(2)$(3) && rm $(2)$(3)
zip = cd $(build_dir) && zip $(2)_$(1)_$(version).zip $(2)$(3) && rm $(2)$(3)

.PHONY: all macos linux windows clean dep fmt test install

default: install

all: macos linux windows

clean:
	rm -rf $(build_dir)/

dep:
	go get -v -d $(pkg_base)/...

fmt: dep
	go fmt $(pkg_base)/...

test: dep
	go test -v $(pkg_base)/...

install: dep
	go install $(pkg_cmd)

##### LINUX BUILDS #####
linux: test build/linux_arm build/linux_arm64 build/linux_amd64

build/linux_amd64:
	$(call build,linux,amd64,$(name_cmd),,$(pkg_cmd))
	$(call tar,linux_amd64,$(name_cmd),)

build/linux_arm:
	$(call build,linux,arm,$(name_cmd),,$(pkg_cmd))
	$(call tar,linux_arm,$(name_cmd),)

build/linux_arm64:
	$(call build,linux,arm64,$(name_cmd),,$(pkg_cmd))
	$(call tar,linux_arm64,$(name_cmd),)

##### MACOS BUILDS #####
macos: test build/macos_amd64 build/macos_arm64

build/macos_amd64:
	$(call build,darwin,amd64,$(name_cmd),,$(pkg_cmd))
	$(call tar,macos_amd64,$(name_cmd),)

build/macos_arm64:
	$(call build,darwin,arm64,$(name_cmd),,$(pkg_cmd))
	$(call tar,macos_arm64,$(name_cmd),)

##### WINDOWS BUILDS #####
windows: test build/windows_amd64

build/windows_amd64:
	$(call build,windows,amd64,$(name_cmd),.exe,$(pkg_cmd))
	$(call zip,windows_amd64,$(name_cmd),.exe)
