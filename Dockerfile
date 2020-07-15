FROM golang:alpine as builder
WORKDIR /go/src/dms
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-s -w -static"' -o /go/bin/dms .

FROM scratch
COPY --from=builder /go/bin/dms /usr/local/bin/dms

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/dms"]
