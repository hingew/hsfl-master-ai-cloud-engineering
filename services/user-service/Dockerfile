#############################
# Builder stage
#############################
FROM golang:1.21-alpine as builder

ENV POSTGRES_HOST=
ENV POSTGRES_PORT=
ENV POSTGRES_USERNAME=
ENV POSTGRES_PASSWORD=
ENV POSTGRES_DBNAME=
ENV AUTH_PRIVATE_KEY=
ENV AUTH_PUBLIC_KEY=

WORKDIR /build
COPY ./lib ./lib
COPY ./services/user-service ./services/user-service

WORKDIR /build/services/user-service
RUN go build -o /app/user-service

#############################
# Runtime stage
#############################
FROM alpine:3.18 as runtime

RUN mkdir -p /app/
WORKDIR /app

COPY --from=builder /app/ /app
EXPOSE 3000
EXPOSE 3001

ENTRYPOINT /app/user-service
