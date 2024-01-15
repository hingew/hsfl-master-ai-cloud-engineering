FROM golang:1.21-alpine
RUN apk add --update nodejs npm
RUN npm i -g yarn

ENV WEB_ENDPOINT=

WORKDIR /app
COPY ./lib ./lib
COPY ./services/web-service ./services/web-service

WORKDIR /app/services/web-service
RUN yarn
RUN yarn generate
RUN yarn build
RUN go mod tidy
RUN go build -o ./main

EXPOSE 3000
CMD ["/app/services/web-service/main"]
