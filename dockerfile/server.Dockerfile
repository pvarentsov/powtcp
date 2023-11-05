FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum /
RUN go mod download

COPY cmd/server/ /cmd/server/
COPY internal/ /internal/
COPY config/demo.yaml /config/

RUN go build -o /server /cmd/server/*.go

EXPOSE 8080

CMD ["/server", "--config", "/config/demo.yaml"]