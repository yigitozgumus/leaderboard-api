FROM golang:1.16.3-alpine3.13
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build ./cmd/webserver
CMD ["/app/webserver"]
