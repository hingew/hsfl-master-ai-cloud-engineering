#############################
# Builder stage
#############################
FROM golang:1.21-alpine as builder

ENV TEMPLATE_GRPC_ENDPOINT=

RUN apk --no-cache --update-cache --available upgrade \
    && apk add terminus-font ttf-inconsolata ttf-dejavu font-noto ttf-font-awesome


WORKDIR /build
COPY ./lib ./lib
COPY ./services/creation-service ./services/creation-service

WORKDIR /build/services/creation-service
RUN go build -o /app/creation-service

#############################
# Runtime stage
#############################
FROM alpine:3.18 as runtime

RUN apk --no-cache --update-cache --available upgrade \
    && apk add terminus-font ttf-inconsolata ttf-dejavu font-noto ttf-font-awesome

RUN mkdir -p /app/
WORKDIR /app

COPY --from=builder /app/ /app
EXPOSE 3000

ENTRYPOINT /app/creation-service
