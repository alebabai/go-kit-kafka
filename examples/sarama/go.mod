module github.com/alebabai/go-kit-kafka/examples/sarama

go 1.16

require (
	github.com/Shopify/sarama v1.32.0
	github.com/alebabai/go-kit-kafka v0.0.0
	github.com/alebabai/go-kit-kafka/examples/common v0.0.0
	github.com/go-kit/kit v0.13.0
	github.com/go-kit/log v0.2.1
)

replace github.com/alebabai/go-kit-kafka => ../..

replace github.com/alebabai/go-kit-kafka/examples/common => ../common
