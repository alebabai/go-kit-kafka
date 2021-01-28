# go-kit-kafka

> Apache Kafka integration module for go-kit

[![build](https://img.shields.io/github/workflow/status/alebabai/go-kit-kafka/CI)](https://github.com/alebabai/go-kit-kafka/actions?query=workflow%3ACI)
[![version](https://img.shields.io/github/go-mod/go-version/alebabai/go-kit-kafka)](https://golang.org/)
[![report](https://goreportcard.com/badge/github.com/alebabai/go-kit-kafka)](https://goreportcard.com/report/github.com/alebabai/go-kit-kafka)
[![coverage](https://img.shields.io/codecov/c/github/alebabai/go-kit-kafka)](https://codecov.io/github/alebabai/go-kit-kafka)
[![release](https://img.shields.io/github/release/alebabai/go-kit-kafka.svg)](https://github.com/alebabai/go-kit-kafka/releases)
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

## Reference

TBD

## Usage

### Consumer

To use consumer transport abstractions the following adapters for the chosen Apache Kafka
client library should be implemented:

```go
type Message interface {
	Topic() string
	Partition() int
	Offset() int64
	Key() []byte
	Value() []byte
	Headers() []Header
	Timestamp() time.Time
}


type Header interface {
	Key() string
	Value() []byte
}


type Reader interface {
	ReadMessage(ctx context.Context) (Message, error)
	Committer
	io.Closer
}

type Committer interface {
    Commit(ctx context.Context) error
}

```

## Examples

Go to [Examples](examples).
