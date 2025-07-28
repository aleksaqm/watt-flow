FROM golang:latest

RUN apt-get update && apt-get install -y wkhtmltopdf

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /go-gin-api

EXPOSE 5000

CMD ["/go-gin-api"]
