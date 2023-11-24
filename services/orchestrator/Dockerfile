FROM golang:1.21-alpine

ENV BALANCER_ALGORITHM=

WORKDIR /build
COPY ./lib ./lib
COPY ./services/orchestrator ./services/orchestrator
WORKDIR /build/services/orchestrator
RUN go build -o /app/orchestrator
EXPOSE 3000

CMD ["/app/orchestrator"]
