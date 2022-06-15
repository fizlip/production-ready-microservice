# production-ready-microservice

This microservice is built using the <a href="https://github.com/go-kit">go-kit</a>
microservice toolkit.

The simple sevice will convert a string into a hash using SHA256.

# Usage

```
go run .
```

```
curl -d '{"id":"test"}' -X POST -H 'Content-Type: application/json' http://localhost:8080/hash/

{"hash":"n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg="}
```
