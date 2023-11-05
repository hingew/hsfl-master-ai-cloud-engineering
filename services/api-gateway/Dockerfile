FROM golang:1.21-alpine

WORKDIR /app
COPY ./src/api-gateway ./src/api-gateway

WORKDIR /app/src/api-gateway
RUN yarn
RUN yarn build
RUN go build -o ./main

EXPOSE 3000
CMD ["/app/src/api-gateway/main"]