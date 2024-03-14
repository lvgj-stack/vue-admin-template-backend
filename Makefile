.DEFAULT_GOAL := all

.PHONY: all

all: tidy format build

ROOT_PACKAGE=github.com/Mr-LvGJ/gobase

include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/golang.mk
include scripts/make-rules/tools.mk


.PHONY: tidy
tidy:
	@$(GO) mod tidy

.PHONY: format
format: tools.verify.golines tools.verify.goimports
	@echo "=======> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=120 --reformat-tags --shorten-comments --ignore-generated .
	@$(GO) mod edit -fmt

.PHONY: build
build:
	@$(MAKE) go.build