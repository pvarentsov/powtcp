FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum /
RUN go mod download

COPY cmd/client/ /cmd/client/
COPY internal/ /internal/
COPY config/ /config/

RUN go build -o /client /cmd/client/*.go

CMD ["/client", "--config", "/config/demo.yaml"]