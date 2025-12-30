# Build.
FROM golang:1.23-alpine AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add build-base
RUN go mod download
COPY . /app
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/entrypoint
EXPOSE 3000
ENTRYPOINT ["/app/entrypoint"]
