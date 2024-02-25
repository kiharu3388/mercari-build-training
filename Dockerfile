FROM golang:1.22-alpine

RUN apk add --no-cache gcc musl-dev

RUN addgroup -S mercari && adduser -S trainee -G mercari
# RUN chown -R trainee:mercari /path/to/db

WORKDIR /app

COPY go /app/go
COPY db /app/db

RUN chown -R trainee:mercari /app/db

WORKDIR /app/go

RUN go mod tidy

RUN CGO_ENABLED=1 go build -o ./mercari-build-training ./app/*.go

CMD go run app/main.go
