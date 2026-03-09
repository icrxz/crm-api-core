FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server ./cmd/api

FROM scratch AS prod
COPY --from=builder /app/server /server
COPY --from=builder /app/resources /resources
COPY --from=builder /app/migrations /app/migrations

ENTRYPOINT ["/server"]

FROM builder AS dev

RUN go install github.com/air-verse/air@latest
