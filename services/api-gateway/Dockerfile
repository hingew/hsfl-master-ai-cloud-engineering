#############################
# Builder stage
#############################
FROM golang:1.21-alpine as builder

ENV TEMPLATE_ENDPOINT=
ENV WEB_ENDPOINT=
ENV CREATION_ENDPOINT=
ENV USER_ENDPOINT=

WORKDIR /build
COPY ./services/api-gateway ./services/api-gateway

WORKDIR /build/services/api-gateway
RUN go mod tidy
RUN go build -o /app/api-gateway

#############################
# Runtime stage
#############################
FROM alpine:3.18 as runtime

RUN mkdir -p /app/
WORKDIR /app

COPY --from=builder /app/ /app
EXPOSE 3000

ENTRYPOINT /app/api-gateway

