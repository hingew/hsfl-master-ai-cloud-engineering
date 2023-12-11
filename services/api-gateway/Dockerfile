FROM golang:1.21-alpine

ENV TEMPLATE_ENDPOINT=
ENV WEB_ENDPOINT=
ENV CREATION_ENDPOINT=
ENV USER_ENDPOINT=
ENV REQUEST_ENDPOINT=

WORKDIR /build
COPY ./services/api-gateway ./services/api-gateway

WORKDIR /build/services/api-gateway
RUN go mod tidy
RUN go build -o /app/api-gateway

EXPOSE 3000
CMD ["/app/api-gateway"]
