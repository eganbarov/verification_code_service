# Stage 1: Build the Go application
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /go-app .

# Stage 2: Run the application in a minimal image
FROM alpine:latest
COPY --from=builder /go-app /go-app
EXPOSE 8080
CMD [ "/go-app" ]