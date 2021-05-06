# Compile stage
FROM golang:1.16 AS build-env
#FROM golang:1.16-alpine AS build-env

#RUN apk add --no-cache build-base
ADD . /go/src/application
WORKDIR /go/src/application
RUN go build -race -o /application -mod=vendor
#RUN go build -o /application -mod=vendor

#CMD go run -race main.go

# Final stage
#FROM alpine:3.11
#
#WORKDIR /
#COPY --from=build-env /application /

CMD ["/application"]
