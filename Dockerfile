FROM golang:alpine as builder

RUN apk --update add ca-certificates

WORKDIR /go/src/dms

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-s -w -static"' -o /go/bin/dms .

FROM scratch

# Add in certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Add the binary
COPY --from=builder /go/bin/dms /usr/local/bin/dms

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/dms"]
