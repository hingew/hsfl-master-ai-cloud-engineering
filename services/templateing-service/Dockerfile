FROM golang:1.21-alpine

ENV TEMPLATE_ENDPOINT=

WORKDIR /app
COPY ./lib ./lib
COPY ./services/templateing-service ./services/templateing-service

WORKDIR /app/services/templateing-service
RUN go build -o ./main

EXPOSE 3000
CMD ["/app/services/templateing-service/main"]