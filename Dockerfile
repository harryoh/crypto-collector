FROM golang:1.14-alpine AS builder

RUN apk --no-cache add ca-certificates

WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -v

#####

FROM scratch

WORKDIR /opt/app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/src/app/crypto-collector .
COPY ui/build ./ui/build

EXPOSE 8080

CMD ["./crypto-collector"]
