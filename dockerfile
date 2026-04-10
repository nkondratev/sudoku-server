FROM golang:1.26.1 AS builder
WORKDIR /build
COPY . .
RUN go mod tidy && go build -o /main .

FROM alpine:3.23
COPY --from=builder main /bin/main
EXPOSE 8080
ENTRYPOINT ["/bin/main"]
