FROM golang:1.21-alpine

ENV BALANCER_ALGORITHM=

WORKDIR /build
COPY ./lib ./lib
COPY ./services/orchestrated-templateing-service ./services/orchestrated-templateing-service

WORKDIR /build/services/orchestrated-templateing-service
RUN go build -o /app/orchestrated-templateing-service

EXPOSE 3000
CMD ["/app/orchestrated-templateing-service"]
