# Build stage
FROM golang:alpine AS build-env

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# Run stage
FROM alpine:latest AS final

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build-env /app/main .

CMD ["./main"]