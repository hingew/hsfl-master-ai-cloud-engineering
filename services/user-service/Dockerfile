FROM golang:1.21-alpine

WORKDIR /build
COPY ./lib ./lib
COPY ./services/user-service ./services/user-service

WORKDIR /build/services/user-service
RUN go build -o /app/user-service

# Add a empty files to make single file volumes work
WORKDIR /app
RUN touch config.yml cert.pem

EXPOSE 3000
CMD ["/app/user-service"]
