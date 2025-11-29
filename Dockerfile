ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY . .
RUN go mod download && go mod verify
RUN go build -v -o /run-app ./cmd/app

FROM debian:bookworm
WORKDIR /app
COPY --from=builder /run-app /usr/local/bin/run-app
COPY --from=builder /usr/src/app/docs /app/docs
CMD ["run-app"]
