# Build stage
FROM golang:alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY public ./public

EXPOSE 8080
CMD ["/app/main"]
