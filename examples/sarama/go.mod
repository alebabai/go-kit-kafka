module github.com/alebabai/go-kit-kafka/v2/examples/sarama

go 1.22.0

require (
	github.com/IBM/sarama v1.42.2
	github.com/alebabai/go-kafka v0.3.0
	github.com/alebabai/go-kafka/adapter/sarama v0.3.0
	github.com/alebabai/go-kit-kafka/v2/examples/common v0.0.0
	github.com/go-kit/kit v0.13.0
	github.com/go-kit/log v0.2.1
)

require (
	github.com/alebabai/go-kit-kafka/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.6.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
)

replace github.com/alebabai/go-kit-kafka/v2 => ../..

replace github.com/alebabai/go-kit-kafka/v2/examples/common => ../common
