FROM golang:1.21-alpine

WORKDIR /build
COPY ./lib ./lib
COPY ./services/creation-service ./services/creation-service

WORKDIR /build/services/creation-service
RUN go build -o /app/creation-service

# Add a empty config.yml to make single file volumes work
WORKDIR /app
RUN touch config.yml

EXPOSE 3000
CMD ["/app/creation-service"]
