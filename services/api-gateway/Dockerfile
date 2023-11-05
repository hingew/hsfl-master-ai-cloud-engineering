FROM golang:1.21-alpine

ENV TEMPLATE_ENDPOINT=
ENV WEB_ENDPOINT=
ENV CREATION_ENDPOINT=
ENV USER_ENDPOINT=

WORKDIR /app
COPY ./services/api-gateway ./services/api-gateway

WORKDIR /app/services/api-gateway
RUN go mod tidy
RUN go build -o ./main

EXPOSE 3000
CMD ["/app/services/api-gateway/main"]