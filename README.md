# GO OK

A Simple go web server that returns a bit of useful data for testing
Kubernetes nodes.

## Run

```bash
go get github.com/gin-gonic/gin
go get github.com/nu7hatch/gouuid

GIN_MODE=release go run ./ok.go
```

## Build Docker Container

```bash
docker build -t go-ok .
```