FROM golang:1.23.2-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /go-simulator

EXPOSE 5000

CMD ["/go-simulator"]

