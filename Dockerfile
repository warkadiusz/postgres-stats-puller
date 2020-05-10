FROM golang:1.13-alpine

WORKDIR /go/src/ask
COPY . /go/src/ask

RUN ls && \
    apk add git && \
    go get github.com/joho/godotenv && \
    go get github.com/lib/pq && \
    go get github.com/influxdata/influxdb1-client && \
    go get github.com/influxdata/influxdb1-client/v2 && \
    go get github.com/gorilla/mux && \
    go get github.com/gorilla/websocket


CMD ["go", "run", "server.go"]