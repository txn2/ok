FROM golang:latest AS builder

RUN mkdir -p /go/src/github.com/txn2/ok
COPY . /go/src/github.com/txn2/ok

RUN go get ...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /go/bin/ok ./src/github.com/txn2/ok

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/ok /ok

WORKDIR /

ENTRYPOINT ["/ok"]
