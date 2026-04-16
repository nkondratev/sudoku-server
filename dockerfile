FROM golang:1.26.1 AS builder
WORKDIR /build
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /main .

FROM alpine:3.23
COPY --from=builder /main /app/main

EXPOSE 8080
ENTRYPOINT ["/app/main"]
