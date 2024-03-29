FROM golang:1.14-alpine AS builder

RUN apk --no-cache add ca-certificates

WORKDIR /usr/src/app
COPY . .
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s' -v

#####

FROM alpine

WORKDIR /opt/app

RUN apk add --no-cache tzdata

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/src/app/crypto-collector .
# COPY ui/dist ./ui/dist

EXPOSE 8080

CMD ["./crypto-collector"]
