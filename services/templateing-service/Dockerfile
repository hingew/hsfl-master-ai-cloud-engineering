FROM golang:1.21-alpine

ENV TEMPLATE_ENDPOINT=

WORKDIR /build
COPY ./lib ./lib
COPY ./services/templateing-service ./services/templateing-service

WORKDIR /build/services/templateing-service
RUN go build -o /app/templateing-service

EXPOSE 3000
CMD ["/app/templateing-service"]
