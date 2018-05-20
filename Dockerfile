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

CMD ["/app"]
