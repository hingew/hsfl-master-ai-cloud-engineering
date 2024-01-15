FROM golang:1.21-alpine

RUN apk --no-cache --update-cache --available upgrade \
    && apk add git bash terminus-font ttf-inconsolata ttf-dejavu font-noto ttf-font-awesome


WORKDIR /build
COPY ./lib ./lib
COPY ./services/creation-service ./services/creation-service

WORKDIR /build/services/creation-service
RUN go build -o /app/creation-service

WORKDIR /app

EXPOSE 3000
CMD ["/app/creation-service"]
