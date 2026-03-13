FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download || go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/tamago ./tamago

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/tamago .
COPY DUNGEON.md DUNGEON.md

ENV TAMAGO_CHILD=1

VOLUME ["/app"]

ENTRYPOINT ["./tamago"]
