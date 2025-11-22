FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /ase-to-http-api

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /ase-to-http-api /ase-to-http-api

EXPOSE 80

CMD ["/ase-to-http-api"]
