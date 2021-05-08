FROM golang:1.16-alpine

ADD . /go/src/application
WORKDIR /go/src/application
RUN mkdir -p /usr/local/content && go build -o /application -mod=vendor

CMD ["/application"]
