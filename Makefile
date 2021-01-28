MODE						:= local

GO							:= go
GO_VERSION					:= 1.15

GOLANGCI_LINT				:= golangci-lint
GOLANGCI_LINT_VERSION		:= 1.29.0

PACKAGES					:= ./...
GO_COVER_PROFILE			:= coverage.out

ifeq ($(MODE),docker)
	GO_DOCKER_IMAGE 				:= library/golang:$(GO_VERSION)
	GO 								:= docker run --rm -v $(CURDIR):/app -v $(GOPATH)/pkg/mod:/go/pkg/mod -w /app $(GO_DOCKER_IMAGE) go

	GOLANGCI_LINT_DOCKER_IMAGE		:= golangci/golangci-lint:v${GOLANGCI_LINT_VERSION}
	GOLANGCI_LINT					:= docker run --rm -v $(CURDIR):/app -w /app $(GOLANGCI_LINT_DOCKER_IMAGE) golangci-lint run -v
endif

.PHONY: all
all: build test

.PNONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: deps
deps:
	$(GO) mod tidy -v

.PHONY: prepare
prepare: deps fmt

.PHONY: build
build: prepare
	$(GO) build -v $(PACKAGES)

.PHONY: install
install: prepare
	$(GO) install -v $(PACKAGES)

.PHONY: test
test: prepare
	$(GO) test -v -race -coverprofile=$(GO_COVER_PROFILE) $(PACKAGES)

.PHONY: coverage
coverage: test
	$(GO) tool cover -func=$(GO_COVER_PROFILE) -o coverage.txt
	$(GO) tool cover -html=$(GO_COVER_PROFILE) -o coverage.html

.PHONY: lint
lint: prepare
	$(GOLANGCI_LINT) run -v