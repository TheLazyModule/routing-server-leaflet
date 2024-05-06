# Golang build stage
FROM golang:alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Final stage
FROM alpine:3.19
WORKDIR /app
# Copy built Go application and migrate tool from the builder stage
COPY --from=builder /app/main /app/main
COPY public  /app/public
COPY app.yaml /app/
COPY app.env /app/
#COPY wait-for.sh  /app/
#COPY startup.sh /app/
#RUN chmod +x /app/wait-for.sh
RUN chmod +x /app/main
EXPOSE 8080
ENTRYPOINT ["/app/main"]
