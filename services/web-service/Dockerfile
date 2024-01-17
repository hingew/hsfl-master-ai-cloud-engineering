#############################
# Builder stage
#############################
FROM golang:1.21-alpine as builder
RUN apk add --update nodejs npm
RUN npm i -g yarn

WORKDIR /app
COPY ./lib ./lib
COPY ./services/web-service ./services/web-service

WORKDIR /app/services/web-service
RUN yarn install --frozen-lockfile --no-cache
RUN yarn generate
RUN yarn build

RUN go mod tidy
RUN go build -o ./web-service

#############################
# Runtime stage
#############################
FROM alpine:3.18 as runtime

RUN mkdir -p /app/public
WORKDIR /app

COPY --from=builder /app/services/web-service/web-service /app
COPY --from=builder /app/services/web-service/public/ /app/public
EXPOSE 3000

ENTRYPOINT /app/web-service

