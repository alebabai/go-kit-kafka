GO							?= @go

PACKAGES					?= ./...
COVER_PROFILE				?= coverage.out

GOLANGCI_LINT				?= @golangci-lint

.PHONY: all
all: build test

.PNONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: mod/download
mod/download:
	$(GO) mod download

.PHONY: prepare
prepare: mod/download fmt

.PHONY: build
build: prepare
	$(GO) build -v $(PACKAGES)

.PHONY: install
install: prepare
	$(GO) install -v $(PACKAGES)

.PHONY: test
test: prepare
	$(GO) test -v -race -coverprofile=$(COVER_PROFILE) $(PACKAGES)

.PHONY: coverage
coverage: test
	$(GO) tool cover -func=$(COVER_PROFILE) -o coverage.txt
	$(GO) tool cover -html=$(COVER_PROFILE) -o coverage.html

.PHONY: lint
lint: prepare
	$(GOLANGCI_LINT) run -v
