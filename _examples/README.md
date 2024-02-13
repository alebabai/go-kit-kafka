# Examples

> Examples of using go-kit-kafka with various Apache Kafka client libraries

## Implementations

- [confluent](confluent)
- [sarama](sarama)

## Reference

### Services

- _producer_ - the service that provide an endpoint to produce an event

- _consumer_ - the service that able to consume events from kafka topic, store them the in inmemory storage and provide
  an endpoint to list all consumed events

## Usage

### Docker

Bootstrap full project using docker compose:

```bash
docker compose -f <example>/compose.yml up
```

To produce an event send the following request:

```bash
curl -X POST http://localhost:8080/events

```

To view all events send the following request:

```bash
curl -X GET http://localhost:8081/events

```
