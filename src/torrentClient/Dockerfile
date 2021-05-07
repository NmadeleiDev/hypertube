#FROM golang:1.16
FROM golang:1.16-alpine

ADD . /go/src/application
WORKDIR /go/src/application
#RUN go build -race -o /application -mod=vendor
RUN go build -o /application -mod=vendor

CMD ["/application"]
