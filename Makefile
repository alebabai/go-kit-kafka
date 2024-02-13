GO							?= @go

PACKAGES					?= ./...
TEST_PACKAGES				?= ./...
COVER_PACKAGES				?= $(shell echo $(TEST_PACKAGES) | tr " " ",")
COVER_PROFILE				?= coverage.out

GOLANGCI_LINT				?= @golangci-lint

.PHONY: all
all: build test

.PNONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: mod/tidy
mod/tidy:
	$(GO) mod tidy

.PHONY: prepare
prepare: mod/tidy fmt

.PHONY: build
build: prepare
	$(GO) build -v $(PACKAGES)

.PHONY: install
install: prepare
	$(GO) install -v $(PACKAGES)

.PHONY: test
test: prepare
	$(GO) test -v -race -coverpkg $(COVER_PACKAGES) -coverprofile $(COVER_PROFILE) $(TEST_PACKAGES)

.PHONY: coverage
coverage: test
	$(GO) tool cover -func $(COVER_PROFILE) -o coverage.txt
	$(GO) tool cover -html $(COVER_PROFILE) -o coverage.html

.PHONY: lint
lint: prepare
	$(GOLANGCI_LINT) run -v
