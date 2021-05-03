# go-kit-kafka

> Apache Kafka integration module for go-kit

[![build](https://img.shields.io/github/workflow/status/alebabai/go-kit-kafka/CI)](https://github.com/alebabai/go-kit-kafka/actions?query=workflow%3ACI)
[![version](https://img.shields.io/github/go-mod/go-version/alebabai/go-kit-kafka)](https://golang.org/)
[![report](https://goreportcard.com/badge/github.com/alebabai/go-kit-kafka)](https://goreportcard.com/report/github.com/alebabai/go-kit-kafka)
[![coverage](https://img.shields.io/codecov/c/github/alebabai/go-kit-kafka)](https://codecov.io/github/alebabai/go-kit-kafka)
[![tag](https://img.shields.io/github/tag/alebabai/go-kit-kafka.svg)](https://github.com/alebabai/go-kit-kafka/tags)
[![reference](https://pkg.go.dev/badge/github.com/alebabai/go-kit-kafka.svg)](https://pkg.go.dev/github.com/alebabai/go-kit-kafka)

## Getting started

Go modules are supported.  

Manual install:

```bash
go get -u github.com/alebabai/go-kit-kafka
```

Golang import:

```go
import "github.com/alebabai/go-kit-kafka/kafka"
```

## Usage

To use consumer/producer transport abstractions the following adapters for the chosen Apache Kafka
client library should be implemented:

```go
type Message struct {
    Topic     string
    Partition int32
    Offset    int64
    Key       []byte
    Value     []byte
    Headers   []Header
    Timestamp time.Time
}

type Header struct {
    Key   []byte
    Value []byte
}
```

## Examples

Go to [Examples](examples).
