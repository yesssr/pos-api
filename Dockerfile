ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY . .

# install dependencies
RUN go mod download && go mod verify

# build executable dari cmd/app
RUN go build -v -o /run-app ./cmd/app

FROM debian:bookworm
COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/docs /app/docs
CMD ["run-app"]
