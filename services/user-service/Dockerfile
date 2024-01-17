FROM golang:1.21-alpine

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

WORKDIR /app

EXPOSE 3000
EXPOSE 3001
CMD ["/app/user-service"]
