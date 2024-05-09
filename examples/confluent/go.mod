module github.com/alebabai/go-kit-kafka/v2/examples/confluent

go 1.22

require (
	github.com/alebabai/go-kafka v0.3.2
	github.com/alebabai/go-kafka/adapter/confluent v0.3.2
	github.com/alebabai/go-kit-kafka/v2/examples/common v0.0.0
	github.com/confluentinc/confluent-kafka-go/v2 v2.4.0
	github.com/go-kit/kit v0.13.0
	github.com/go-kit/log v0.2.1
)

require (
	github.com/alebabai/go-kit-kafka/v2 v2.0.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
)

replace github.com/alebabai/go-kit-kafka/v2 => ../..

replace github.com/alebabai/go-kit-kafka/v2/examples/common => ../common
