module github.com/alebabai/go-kit-kafka/v2/examples/sarama

go 1.22

toolchain go1.22.0

require (
	github.com/Shopify/sarama v1.32.0
	github.com/alebabai/go-kafka v0.3.0
	github.com/alebabai/go-kit-kafka/v2 v2.0.0
	github.com/alebabai/go-kit-kafka/v2/examples/common v0.0.0
	github.com/go-kit/kit v0.13.0
	github.com/go-kit/log v0.2.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.14.4 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/alebabai/go-kit-kafka/v2 => ../..

replace github.com/alebabai/go-kit-kafka/v2/examples/common => ../common
