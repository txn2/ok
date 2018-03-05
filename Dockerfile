FROM golang:latest AS builder

RUN mkdir -p /go/src/github.com/cjimti/go-ok
COPY . /go/src/github.com/cjimti/go-ok

RUN go get ...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /go/bin/ok ./src/github.com/cjimti/go-ok

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/ok /ok

WORKDIR /

ENTRYPOINT ["/ok"]
