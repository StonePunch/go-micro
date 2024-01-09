FROM golang:1.21.3-alpine3.18 as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# Probably not needed, turns results of the build into an executable
RUN chmod +x /app/brokerApp

# Build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

# Run executable created in the first image and then transfered to the second, smaller image
CMD [ "/app/brokerApp" ]