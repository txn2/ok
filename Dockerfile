FROM golang:1.15.2-alpine3.12 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir -p /go/src/ \
 && mkdir -p /go/bin \
 && mkdir -p /go/pkg

ENV PATH=/go/bin:$PATH

RUN mkdir -p /go/src/ok/
ADD . /go/src/ok/

WORKDIR /go/src/ok/

RUN go build -ldflags "-extldflags \"-static\"" -o /go/bin/ok ok.go

FROM alpine:3.12

COPY --from=builder /go/bin/ok /ok

WORKDIR /

ENTRYPOINT ["/ok"]
