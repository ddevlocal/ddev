# Makefile for a standard golang repo with associated container

# Circleci doesn't seem to provide a decent way to add to path, just adding here, for case where
# linux build and linuxbrew is installed.
export PATH := $(EXTRA_PATH):$(PATH)

BUILD_BASE_DIR ?= $(PWD)

GOTMP=.gotmp
SHELL = /bin/bash
PWD = $(shell pwd)
GOFILES = $(shell find $(SRC_DIRS) -type f)
.PHONY: darwin_amd64 darwin_arm64 darwin_amd64_notarized darwin_arm64_notarized darwin_arm64_signed darwin_amd64_signed linux_amd64 linux_arm64 linux_arm windows_amd64 windows_arm64

# Expands SRC_DIRS into the common golang ./dir/... format for "all below"
SRC_AND_UNDER = $(patsubst %,./%/...,$(SRC_DIRS))

GOLANGCI_LINT_ARGS ?= --out-format=line-number --disable-all --enable=gofmt --enable=govet --enable=revive --enable=errcheck --enable=staticcheck --enable=ineffassign --enable=varcheck --enable=deadcode

WINDOWS_GSUDO_VERSION=v0.7.3
WINNFSD_VERSION=2.4.0
NSSM_VERSION=2.24-101-g897c7ad
MKCERT_VERSION=v1.4.6

TESTTMP=/tmp/testresults

# This repo's root import path (under GOPATH).
PKG := github.com/drud/ddev

# Top-level directories to build
SRC_DIRS := cmd pkg


VERSION := $(shell git describe --tags --always --dirty)
# Some things insist on having the version without the leading 'v', so provide a
# $(NO_V_VERSION) without it.
# no_v_version removes the front v, for Chocolatey mostly
NO_V_VERSION=$(shell echo $(VERSION) | awk -F"-" '{ OFS="-"; sub(/^./, "", $$1); printf $$0; }')
GITHUB_ORG := drud

BUILD_OS = $(shell go env GOHOSTOS)
BUILD_ARCH = $(shell go env GOHOSTARCH)
DEFAULT_BUILD=$(shell go env GOHOSTOS)_$(shell go env GOHOSTARCH)

build:
	goreleaser build --snapshot --skip-validate --rm-dist --single-target

all:
	goreleaser build --snapshot --rm-dist

release:
	goreleaser build --rm-dist


#todo
# - build all if not specific singletarget, bhow?
# - only skip validate on dirty
# - Build other things like the completions generator and completions themselves
# - Build windows installer (separate build I guess)

TEST_TIMEOUT=4h
BUILD_ARCH = $(shell go env GOARCH)

DDEVNAME=ddev
SHASUM=shasum -a 256
ifeq ($(BUILD_OS),windows)
	DDEVNAME=ddev.exe
	SHASUM=sha256sum
endif

DDEV_PATH=$(PWD)/$(GOTMP)/bin/$(BUILD_OS)_$(BUILD_ARCH)
DDEV_BINARY_FULLPATH=$(DDEV_PATH)/$(DDEVNAME)

# Override test section with tests specific to ddev
test: testpkg testcmd

testcmd: $(DEFAULT_BUILD) setup
	@echo LDFLAGS=$(LDFLAGS)
	@echo DDEV_BINARY_FULLPATH=$(DDEV_BINARY_FULLPATH)
	export PATH="$(DDEV_PATH):$$PATH" DDEV_NO_INSTRUMENTATION=true CGO_ENABLED=0 DDEV_BINARY_FULLPATH=$(DDEV_BINARY_FULLPATH); go test $(USEMODVENDOR) -p 1 -timeout $(TEST_TIMEOUT) -v -installsuffix static -ldflags " $(LDFLAGS) " ./cmd/... $(TESTARGS)

testpkg: testnotddevapp testddevapp

testddevapp: $(DEFAULT_BUILD) setup
	export PATH="$(DDEV_PATH):$$PATH" DDEV_NO_INSTRUMENTATION=true CGO_ENABLED=0 DDEV_BINARY_FULLPATH=$(DDEV_BINARY_FULLPATH); go test $(USEMODVENDOR) -p 1 -timeout $(TEST_TIMEOUT) -v -installsuffix static -ldflags " $(LDFLAGS) " ./pkg/ddevapp $(TESTARGS)

testnotddevapp: $(DEFAULT_BUILD) setup
	export PATH="$(DDEV_PATH):$$PATH" DDEV_NO_INSTRUMENTATION=true CGO_ENABLED=0 DDEV_BINARY_FULLPATH=$(DDEV_BINARY_FULLPATH); go test $(USEMODVENDOR) -p 1 -timeout $(TEST_TIMEOUT) -v -installsuffix static -ldflags " $(LDFLAGS) " $(shell find ./pkg -maxdepth 1 -type d ! -name ddevapp ! -name pkg) $(TESTARGS)

testfullsitesetup: $(DEFAULT_BUILD) setup
	export PATH="$(DDEV_PATH):$$PATH" DDEV_NO_INSTRUMENTATION=true CGO_ENABLED=0 DDEV_BINARY_FULLPATH=$(DDEV_BINARY_FULLPATH); go test $(USEMODVENDOR) -p 1 -timeout $(TEST_TIMEOUT) -v -installsuffix static -ldflags " $(LDFLAGS) " ./pkg/ddevapp -run TestDdevFullSiteSetup $(TESTARGS)

setup:
	@mkdir -p $(GOTMP)/{src,pkg/mod/cache,.cache}
	@mkdir -p $(TESTTMP)

# Required static analysis targets used in circleci - these cause fail if they don't work
staticrequired: setup golangci-lint markdownlint mkdocs pyspelling

# Best to install markdownlint-cli locally with "npm install -g markdownlint-cli"
markdownlint:
	@echo "markdownlint: "
	@CMD="markdownlint *.md docs/*.md docs/users 2>&1"; \
	set -eu -o pipefail; \
	if command -v markdownlint >/dev/null 2>&1 ; then \
		$$CMD; \
	else \
		echo "Skipping markdownlint as not installed"; \
	fi

# Best to install mkdocs locally with "sudo pip3 install -r requirements.txt"
mkdocs:
	@echo "mkdocs: "
	@CMD="mkdocs -q build -d /tmp/mkdocsbuild"; \
	if command -v mkdocs >/dev/null 2>&1; then \
		$$CMD ; \
	else \
		echo "Not running mkdocs because it's not installed"; \
	fi

# To see what the docs build will be, you can use `make mkdocs-serve`
# It works best with mkdocs installed, `pip3 install mkdocs`,
# see https://www.mkdocs.org/user-guide/installation/
# But it will also work using docker.
mkdocs-serve:
	if command -v mkdocs >/dev/null ; then \
  		mkdocs serve; \
	else \
		docker run -it -p 8000:8000 -v "$${PWD}:/docs" -e "ADD_MODULES=mkdocs-material mdx_truly_sane_lists mkdocs-git-revision-date-localized-plugin" -e "LIVE_RELOAD_SUPPORT=true"  -e "FAST_MODE=true" -e "DOCS_DIRECTORY=./docs" polinux/mkdocs; \
	fi

# Install markdown-link-check locally with "npm install -g markdown-link-check"
markdown-link-check:
	@echo "markdown-link-check: "
	if command -v markdown-link-check >/dev/null 2>&1; then \
  		find docs *.md -name "*.md" -exec markdown-link-check -q -c markdown-link-check.json {} \; 2>&1  ;\
	else \
		echo "Not running markdown-link-check because it's not installed"; \
	fi

# Best to install pyspelling locally with "sudo -H pip3 install pyspelling pymdown-extensions". Also requries aspell, `sudo apt-get install aspell"
pyspelling:
	@echo "pyspelling: "
	@CMD="pyspelling --config .spellcheck.yml"; \
	set -eu -o pipefail; \
	if command -v pyspelling >/dev/null 2>&1 ; then \
		$$CMD; \
	else \
		echo "Not running pyspelling because it's not installed"; \
	fi

darwin_amd64_signed: $(GOTMP)/bin/darwin_amd64/ddev
	@if [ -z "$(DDEV_MACOS_SIGNING_PASSWORD)" ] ; then echo "Skipping signing ddev for macOS, no DDEV_MACOS_SIGNING_PASSWORD provided"; else echo "Signing $< ..."; \
		set -o errexit -o pipefail; \
		curl -s https://raw.githubusercontent.com/drud/signing_tools/master/macos_sign.sh | bash -s -  --signing-password="$(DDEV_MACOS_SIGNING_PASSWORD)" --cert-file=certfiles/ddev_developer_id_cert.p12 --cert-name="Developer ID Application: Localdev Foundation (9HQ298V2BW)" --target-binary="$<" ; \
	fi
darwin_arm64_signed: $(GOTMP)/bin/darwin_arm64/ddev
	@if [ -z "$(DDEV_MACOS_SIGNING_PASSWORD)" ] ; then echo "Skipping signing ddev for macOS, no DDEV_MACOS_SIGNING_PASSWORD provided"; else echo "Signing $< ..."; \
		set -o errexit -o pipefail; \
		codesign --remove-signature "$(GOTMP)/bin/darwin_arm64/ddev" || true; \
		curl -s https://raw.githubusercontent.com/drud/signing_tools/master/macos_sign.sh | bash -s -  --signing-password="$(DDEV_MACOS_SIGNING_PASSWORD)" --cert-file=certfiles/ddev_developer_id_cert.p12 --cert-name="Developer ID Application: Localdev Foundation (9HQ298V2BW)" --target-binary="$<" ; \
	fi

darwin_amd64_notarized: darwin_amd64_signed
	@if [ -z "$(DDEV_MACOS_APP_PASSWORD)" ]; then echo "Skipping notarizing ddev for macOS, no DDEV_MACOS_APP_PASSWORD provided"; else \
		set -o errexit -o pipefail; \
		echo "Notarizing $(GOTMP)/bin/darwin_amd64/ddev ..." ; \
		curl -sSL -f https://raw.githubusercontent.com/drud/signing_tools/master/macos_notarize.sh | bash -s -  --app-specific-password=$(DDEV_MACOS_APP_PASSWORD) --apple-id=notarizer@localdev.foundation --primary-bundle-id=com.ddev.ddev --target-binary="$(GOTMP)/bin/darwin_amd64/ddev" ; \
	fi
darwin_arm64_notarized: darwin_arm64_signed
	@if [ -z "$(DDEV_MACOS_APP_PASSWORD)" ]; then echo "Skipping notarizing ddev for macOS, no DDEV_MACOS_APP_PASSWORD provided"; else \
		set -o errexit -o pipefail; \
		echo "Notarizing $(GOTMP)/bin/darwin_arm64/ddev ..." ; \
		curl -sSL -f https://raw.githubusercontent.com/drud/signing_tools/master/macos_notarize.sh | bash -s - --app-specific-password=$(DDEV_MACOS_APP_PASSWORD) --apple-id=notarizer@localdev.foundation --primary-bundle-id=com.ddev.ddev --target-binary="$(GOTMP)/bin/darwin_arm64/ddev" ; \
	fi

windows_install: $(GOTMP)/bin/windows_amd64/ddev_windows_installer.$(VERSION).exe

$(GOTMP)/bin/windows_amd64/ddev_windows_installer.$(VERSION).exe: $(GOTMP)/bin/windows_amd64/ddev.exe $(GOTMP)/bin/windows_amd64/sudo_license.txt $(GOTMP)/bin/windows_amd64/mkcert.exe $(GOTMP)/bin/windows_amd64/mkcert_license.txt winpkg/ddev.nsi
	ls -l .gotmp/bin/windows_amd64
	@if [ "$(DDEV_WINDOWS_SIGN)" != "true" ] ; then echo "Skipping signing ddev.exe, DDEV_WINDOWS_SIGN not set"; else echo "Signing windows binaries..." && signtool sign ".gotmp/bin/windows_amd64/ddev.exe" ".gotmp/bin/windows_amd64/mkcert.exe" ".gotmp/bin/windows_amd64/ddev_gen_autocomplete.exe"; fi
	@makensis -DVERSION=$(VERSION) winpkg/ddev.nsi  # brew install makensis, apt-get install nsis, or install on Windows
	@if [ "$(DDEV_WINDOWS_SIGN)" != "true" ] ; then echo "Skipping signing ddev_windows_installer, DDEV_WINDOWS_SIGN not set"; else echo "Signing windows installer binary..." && signtool sign "$@"; fi
	$(SHASUM) $@ >$@.sha256.txt

no_v_version:
	@echo $(NO_V_VERSION)

chocolatey: $(GOTMP)/bin/windows_amd64/ddev_windows_installer.$(VERSION).exe
	rm -rf $(GOTMP)/bin/windows_amd64/chocolatey && cp -r winpkg/chocolatey $(GOTMP)/bin/windows_amd64/chocolatey
	perl -pi -e 's/REPLACE_DDEV_VERSION/$(NO_V_VERSION)/g' $(GOTMP)/bin/windows_amd64/chocolatey/*.nuspec
	perl -pi -e 's/REPLACE_DDEV_VERSION/$(VERSION)/g' $(GOTMP)/bin/windows_amd64/chocolatey/tools/*.ps1
	perl -pi -e 's/REPLACE_GITHUB_ORG/$(GITHUB_ORG)/g' $(GOTMP)/bin/windows_amd64/chocolatey/*.nuspec $(GOTMP)/bin/windows_amd64/chocolatey/tools/*.ps1 #GITHUB_ORG is for testing, for example when the binaries are on rfay acct
	perl -pi -e "s/REPLACE_INSTALLER_CHECKSUM/$$(cat $(GOTMP)/bin/windows_amd64/ddev_windows_installer.$(VERSION).exe.sha256.txt | awk '{ print $$1; }')/g" $(GOTMP)/bin/windows_amd64/chocolatey/tools/*
	if [[ "$(NO_V_VERSION)" =~ -g[0-9a-f]+ ]]; then \
		echo "Skipping chocolatey build on interim version"; \
	else \
		docker run --rm -v "/$(PWD)/$(GOTMP)/bin/windows_amd64/chocolatey:/tmp/chocolatey" -w "//tmp/chocolatey" linuturk/mono-choco pack ddev.nuspec; \
		echo "chocolatey package is in $(GOTMP)/bin/windows_amd64/chocolatey"; \
	fi

$(GOTMP)/bin/windows_amd64/mkcert.exe $(GOTMP)/bin/windows_amd64/mkcert_license.txt:
	curl --fail -sSL -o $(GOTMP)/bin/windows_amd64/mkcert.exe https://github.com/drud/mkcert/releases/download/$(MKCERT_VERSION)/mkcert-$(MKCERT_VERSION)-windows-amd64.exe
	curl --fail -sSL -o $(GOTMP)/bin/windows_amd64/mkcert_license.txt -O https://raw.githubusercontent.com/drud/mkcert/master/LICENSE

$(GOTMP)/bin/windows_amd64/sudo_license.txt:
	set -x
	curl --fail -sSL -o "$(GOTMP)/bin/windows_amd64/sudo_license.txt" "https://raw.githubusercontent.com/gerardog/gsudo/master/LICENSE.txt"

# Best to install golangci-lint locally with "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.31.0"
golangci-lint:
	@echo "golangci-lint: "
	@CMD="golangci-lint run $(GOLANGCI_LINT_ARGS) $(SRC_AND_UNDER)"; \
	set -eu -o pipefail; \
	if command -v golangci-lint >/dev/null 2>&1; then \
		$$CMD; \
	else \
		echo "Skipping golanci-lint as not installed"; \
	fi

version:
	@echo VERSION:$(VERSION)

clean:
	@rm -rf artifacts

# print-ANYVAR prints the expanded variable
print-%: ; @echo $* = $($*)
