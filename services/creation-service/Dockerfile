FROM golang:1.21-alpine

WORKDIR /build
COPY ./lib ./lib
COPY ./services/creation-service ./services/creation-service

WORKDIR /build/services/creation-service
RUN go build -o /app/creation-service

WORKDIR /app

EXPOSE 3000
CMD ["/app/creation-service"]
