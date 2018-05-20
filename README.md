# Docker builder pattern with scratch image

A scratch based docker image template for the golang apps.

## Build
```
docker build -t hello .
```

## Run
```
docker run -p 3000:3000 -e GIN_MODE=debug -e HOSTPORT=:3000 --rm -it hello:latest
```

## Stress Test
```
docker run --net=host --rm -it istio/fortio load -qps 1000 http://localhost:3000/v1/quotes/aapl
```

## Scratch based docker image with certs and vendor deps

```
FROM golang:alpine AS builder
ADD ./vendor /go/vendor
ADD . /go/src/app
WORKDIR /go/src/app
RUN apk update && apk add ca-certificates
RUN CGO_ENABLED=0 go build -tags netgo -ldflags '-w'

FROM scratch
COPY --from=builder /go/src/app/app /app
COPY --from=builder /usr/share/ca-certificates /usr/share/ca-certificates 
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENV HOSTPORT=":8080"

ENTRYPOINT ["/app"]
```

