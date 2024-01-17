FROM golang:1.21-alpine

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

EXPOSE 3000
EXPOSE 3001
CMD ["/app/templateing-service"]
