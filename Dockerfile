FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

FROM scratch AS prod
COPY --from=builder /app/server /server
COPY --from=builder /app/resources /resources
COPY --from=builder /app/migrations /migrations

ENTRYPOINT ["/server"]

FROM builder AS dev

RUN go install github.com/cosmtrek/air@latest
