FROM golang:1.21-alpine

ENV WEB_SERVICE_ENDPOINT=

WORKDIR /build
COPY ./services/request-service ./services/request-service

WORKDIR /build/services/request-service
RUN go mod tidy
RUN go build -o /app/request-service

EXPOSE 3000
CMD ["/app/request-service"]
