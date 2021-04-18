FROM golang:1.16 AS builder

RUN go version

COPY . /builder/
WORKDIR /builder/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/task-tracker ./cmd/task-tracker/main.go

FROM debian:bullseye-slim
RUN apt-get update
RUN apt-get install -y ca-certificates && update-ca-certificates

WORKDIR /app/

COPY --from=builder /builder/.bin/task-tracker .
COPY --from=builder /builder/configs configs/

CMD ["./task-tracker"]
