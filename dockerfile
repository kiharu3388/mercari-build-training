FROM golang:1.22-alpine

RUN addgroup -S mercari && adduser -S trainee -G mercari
# RUN chown -R trainee:mercari /path/to/db

WORKDIR /app

COPY go /app/go
COPY db /app/db

WORKDIR /app/go

RUN chown -R trainee:mercari /app/go
USER trainee

RUN go mod tidy
CMD go run app/main.go
