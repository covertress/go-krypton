# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gkr gkr-cross evm all test travis-test-with-coverage xgo clean
.PHONY: gkr-linux gkr-linux-arm gkr-linux-386 gkr-linux-amd64
.PHONY: gkr-darwin gkr-darwin-386 gkr-darwin-amd64
.PHONY: gkr-windows gkr-windows-386 gkr-windows-amd64
.PHONY: gkr-android gkr-android-16 gkr-android-21

GOBIN = build/bin

CROSSDEPS = https://gmplib.org/download/gmp/gmp-6.0.0a.tar.bz2
GO ?= latest

gkr:
	build/env.sh go install -v $(shell build/flags.sh) ./cmd/gkr
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gkr\" to launch gkr."

gkr-cross: gkr-linux gkr-darwin gkr-windows gkr-android
	@echo "Full cross compilation done:"
	@ls -l $(GOBIN)/gkr-*

gkr-linux: xgo gkr-linux-arm gkr-linux-386 gkr-linux-amd64
	@echo "Linux cross compilation done:"
	@ls -l $(GOBIN)/gkr-linux-*

gkr-linux-arm: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --dest=$(GOBIN) --deps=$(CROSSDEPS) --targets=linux/arm -v $(shell build/flags.sh) ./cmd/gkr
	@echo "Linux ARM cross compilation done:"
	@ls -l $(GOBIN)/gkr-linux-* | grep arm

gkr-linux-386: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --dest=$(GOBIN) --deps=$(CROSSDEPS) --targets=linux/386 -v $(shell build/flags.sh) ./cmd/gkr
	@echo "Linux 386 cross compilation done:"
	@ls -l $(GOBIN)/gkr-linux-* | grep 386

gkr-linux-amd64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --dest=$(GOBIN) --deps=$(CROSSDEPS) --targets=linux/amd64 -v $(shell build/flags.sh) ./cmd/gkr
	@echo "Linux amd64 cross compilation done:"
	@ls -l $(GOBIN)/gkr-linux-* | grep amd64

gkr-darwin: xgo gkr-darwin-386 gkr-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -l $(GOBIN)/gkr-darwin-*

gkr-darwin-386: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --dest=$(GOBIN) --deps=$(CROSSDEPS) --targets=darwin/386 -v $(shell build/flags.sh) ./cmd/gkr
	@echo "Darwin 386 cross compilation done:"
	@ls -l $(GOBIN)/gkr-darwin-* | grep 386

gkr-darwin-amd64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --dest=$(GOBIN) --deps=$(CROSSDEPS) --targets=darwin/amd64 -v $(shell build/flags.sh) ./cmd/gkr
	@echo "Darwin amd64 cross compilation done:"
	@ls -l $(GOBIN)/gkr-darwin-* | grep amd64

gkr-windows: xgo gkr-windows-386 gkr-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -l $(GOBIN)/gkr-windows-*

gkr-windows-386: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --dest=$(GOBIN) --deps=$(CROSSDEPS) --targets=windows/386 -v $(shell build/flags.sh) ./cmd/gkr
	@echo "Windows 386 cross compilation done:"
	@ls -l $(GOBIN)/gkr-windows-* | grep 386

gkr-windows-amd64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --dest=$(GOBIN) --deps=$(CROSSDEPS) --targets=windows/amd64 -v $(shell build/flags.sh) ./cmd/gkr
	@echo "Windows amd64 cross compilation done:"
	@ls -l $(GOBIN)/gkr-windows-* | grep amd64

gkr-android: xgo gkr-android-16 gkr-android-21
	@echo "Android cross compilation done:"
	@ls -l $(GOBIN)/gkr-android-*

gkr-android-16: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --dest=$(GOBIN) --deps=$(CROSSDEPS) --targets=android-16/* -v $(shell build/flags.sh) ./cmd/gkr
	@echo "Android 16 cross compilation done:"
	@ls -l $(GOBIN)/gkr-android-16-*

gkr-android-21: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --dest=$(GOBIN) --deps=$(CROSSDEPS) --targets=android-21/* -v $(shell build/flags.sh) ./cmd/gkr
	@echo "Android 21 cross compilation done:"
	@ls -l $(GOBIN)/gkr-android-21-*

evm:
	build/env.sh $(GOROOT)/bin/go install -v $(shell build/flags.sh) ./cmd/evm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/evm to start the evm."

all:
	build/env.sh go install -v $(shell build/flags.sh) ./...

test: all
	build/env.sh go test ./...

travis-test-with-coverage: all
	build/env.sh build/test-global-coverage.sh

xgo:
	build/env.sh go get github.com/karalabe/xgo

clean:
	rm -fr build/_workspace/pkg/ Godeps/_workspace/pkg $(GOBIN)/*
