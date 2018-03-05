# GO OK

A Simple go web server that returns a bit of useful data for testing
Kubernetes nodes.

## Docker Run
```bash
# run version 1
docker run --rm -p 8080:8080 GIN_MODE=release cjimti/go-ok:v1

# run version 2
docker run --rm -p 8080:8080 GIN_MODE=release cjimti/go-ok:v2

```

Browse to http://localhost:8080

```json
{
    "client_ip": "172.17.0.1",
    "count": 3,
    "message": "ok",
    "time": "2018-03-05T08:38:03.936996398Z",
    "uuid_call": "dddb3561-7273-45ee-5f80-7b022d2bf2e9",
    "uuid_instance": "79defbd7-690e-4fc7-5652-354e1662ff7c",
    "version": 2,
    "version_msg": "version 2"
}
```

## Run Source

```bash
go get github.com/gin-gonic/gin
go get github.com/nu7hatch/gouuid

GIN_MODE=release go run ./ok.go
```

## Build Docker Container

```bash
docker build -t go-ok .
```