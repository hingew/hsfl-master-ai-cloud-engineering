#############################
# Builder stage
#############################
FROM golang:1.21-alpine as builder

ENV POSTGRES_HOST=
ENV POSTGRES_PORT=
ENV POSTGRES_USERNAME=
ENV POSTGRES_PASSWORD=
ENV POSTGRES_DBNAME=
ENV USER_GRPC_ENDPOINT=

WORKDIR /build
COPY ./lib ./lib
COPY ./services/templateing-service ./services/templateing-service

WORKDIR /build/services/templateing-service
RUN go build -o /app/templateing-service

#############################
# Runtime stage
#############################
FROM alpine:3.18 as runtime

RUN mkdir -p /app/
WORKDIR /app

COPY --from=builder /app/ /app
EXPOSE 3000
EXPOSE 3001

ENTRYPOINT /app/templateing-service
