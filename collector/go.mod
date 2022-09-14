module github.com/terracegarden/collector

go 1.18

require (
	github.com/confluentinc/confluent-kafka-go v1.9.1
	github.com/gorilla/mux v1.8.0
	github.com/terracegarden/framework v0.0.0-00010101000000-000000000000
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/xid v1.4.0 // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.opentelemetry.io/otel v1.10.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
)

require (
	github.com/Netflix/go-env v0.0.0-20220526054621-78278af1949d
	github.com/go-pg/pg/v10 v10.10.6 // indirect
	github.com/go-pg/zerochecker v0.2.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/mochi-co/mqtt v1.2.3
	github.com/opentracing/opentracing-go v1.2.0
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/vmihailenco/bufpool v0.1.11 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.4 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opentelemetry.io/otel/trace v1.10.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
	mellium.im/sasl v0.2.1 // indirect
)

replace github.com/terracegarden/framework => ../framework
