FROM golang:1.12-alpine as build

WORKDIR /go/src/github.com/max-samoylov/go-examples/crud
ADD . /go/src/github.com/max-samoylov/go-examples/crud

RUN go build -o crud


FROM alpine:3.9

COPY --from=build /go/src/github.com/max-samoylov/go-examples/crud/crud /srv/crud

ENTRYPOINT [ "/srv/crud" ]