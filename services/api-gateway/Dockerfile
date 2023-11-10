FROM golang:1.21-alpine

ENV TEMPLATE_ENDPOINT=
ENV WEB_ENDPOINT=
ENV CREATION_ENDPOINT=
ENV USER_ENDPOINT=

WORKDIR /build
COPY ./services/api-gateway ./services/api-gateway

WORKDIR /build/services/api-gateway
RUN go mod tidy
RUN go build -o /app/api-gateway

# Add a empty config.yml to make single file volumes work
WORKDIR /app
RUN touch config.yml

EXPOSE 3000
CMD ["/app/api-gateway"]
