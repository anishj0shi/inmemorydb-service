# inmemorydb-service

Inmemorydb-service is a sample lightweight application to show how one can
write an http service in Golang and persist the write request into a InMemory
Database. The Key functionality of this service is to collect Latency information
into the database and retrieve with pagination criteria for plotting latencies. 

This trivial Golang Service  expects httpRequest in the following struct
format

```go
type EventResult struct {
	ID int
	EventId string
	E2ELatency int
	EventType string
}
```

### Deployment

#### Local Deployment
To quickly run the application locally, execute the following command:
```shell script
go run ./cmd/main.go
```

#### Kubernetes Deployment
Deploy this app using `ko`
```shell script
ko apply -f config/ -n <KUBERNETES_NAMESPACE>
```

### Test Application
#### Test on Localhost

Writing to the inmemory DB
```shell script
curl -v -d \
'{"ID": 4, "EventID": "d121", "E2ELatency": 121, "EventType": "order.created"}' \
http://localhost:8080/eventResult
```

Reading Data from the InMemory DB
```shell script
curl -v http://localhost:8080/eventResult\?top\=5\&skip\=0
```
#### Test on Kubernetes

Writing to the inmemory DB
```shell script
curl -v -d \
'{"ID": 4, "EventID": "d121", "E2ELatency": 121, "EventType": "order.created"}' \
http://inmemorydb-service.<YOUR_NAMESPACE>:8080/eventResult
```

Reading Data from the InMemory DB
```shell script
curl -v http://inmemorydb-service.<YOUR_NAMESPACE>:8080/eventResult\?top\=5\&skip\=0
```
