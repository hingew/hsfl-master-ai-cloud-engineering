FROM golang:1.21-alpine

ENV TEMPLATE_ENDPOINT=

WORKDIR /build
COPY ./lib ./lib
COPY ./services/templateing-service ./services/templateing-service

WORKDIR /build/services/templateing-service
RUN go build -o /app/templateing-service

# Add a empty config.yml to make single file volumes work
WORKDIR /app
RUN touch config.yml

EXPOSE 3000
CMD ["/app/templateing-service"]
