FROM golang as builder

WORKDIR /go/src/github.com/afeldman/go-position

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /bin
COPY --from=builder /go/src/github.com/afeldman/go-position/app .

EXPOSE 8888

CMD ["./app", "--release"]
