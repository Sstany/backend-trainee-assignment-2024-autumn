FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY ./ ./

RUN apk update && \
    apk add build-base

RUN go mod download

RUN go build -o avito2024 /app/cmd/core/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/avito2024 /app/avito2024

CMD [ "./avito2024" ]