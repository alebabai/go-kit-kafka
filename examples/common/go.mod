module github.com/alebabai/go-kit-kafka/v2/examples/common

go 1.22

require (
	github.com/alebabai/go-kafka v0.3.0
	github.com/alebabai/go-kit-kafka/v2 v2.0.0
	github.com/go-kit/kit v0.13.0
	github.com/go-kit/log v0.2.1
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
)

require github.com/go-logfmt/logfmt v0.5.1 // indirect

replace github.com/alebabai/go-kit-kafka/v2 => ../../
