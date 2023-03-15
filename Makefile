##@ Linters

GO := go

OUTPUT_DIR := out
LOCALBIN := bin

GOLANGCI_LINT=$(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.51.2

YAMLLINT_VERSION ?= 1.28.0

SHELLCHECK=$(LOCALBIN)/shellcheck
SHELLCHECK_VERSION ?= v0.9.0

GO_LINT_CMD = GOFLAGS="$(GOFLAGS)" GOGC=30 GOCACHE=$(GOCACHE) $(GOLANGCI_LINT) run

build:
	$(GO) build -ldflags="-s -w" -trimpath -o out/kid main.go

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

.PHONY: lint-go
lint-go: $(GOLANGCI_LINT) fmt vet ## Checks Go code
	$(GO_LINT_CMD)

$(GOLANGCI_LINT):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) $(GOLANGCI_LINT_VERSION)

.PHONY: shellcheck
shellcheck: $(SHELLCHECK) ## Download shellcheck locally if necessary.
$(SHELLCHECK): $(OUTPUT_DIR)
ifeq (,$(wildcard $(SHELLCHECK)))
ifeq (,$(shell which shellcheck 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(SHELLCHECK)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH | sed -e 's,amd64,x86_64,g') && \
	curl -Lo $(OUTPUT_DIR)/shellcheck.tar.xz https://github.com/koalaman/shellcheck/releases/download/$(SHELLCHECK_VERSION)/shellcheck-$(SHELLCHECK_VERSION).$${OS}.$${ARCH}.tar.xz ;\
	tar --directory $(OUTPUT_DIR) -xvf $(OUTPUT_DIR)/shellcheck.tar.xz ;\
	find $(OUTPUT_DIR) -name shellcheck -exec cp {} $(SHELLCHECK) \; ;\
	chmod +x $(SHELLCHECK) ;\
	}
else
SHELLCHECK = $(shell which shellcheck)
endif
endif

.PHONY: lint-shell
lint-shell: $(SHELLCHECK) ## Check shell scripts
	find . -name vendor -prune -o -name '*.sh' -print | xargs $(SHELLCHECK) -x

.PHONY: lint-shell-fix
lint-shell-fix: $(SHELLCHECK)
	find * -name vendor -prune -o -name '*.sh' -type f -print | xargs -I@ sh -c "$(SHELLCHECK) -f diff @ | git apply"
