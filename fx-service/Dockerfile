FROM golang:1.12-alpine as build

WORKDIR /go/src/github.com/max-samoylov/go-examples/fx-service
COPY . /go/src/github.com/max-samoylov/go-examples/fx-service

RUN go build -ldflags "-s -w" -o fx-service ./app


FROM alpine:3.9

COPY --from=build /go/src/github.com/max-samoylov/go-examples/fx-service/wait.sh /srv/wait.sh
COPY --from=build /go/src/github.com/max-samoylov/go-examples/fx-service/app/sql/queries.sql /sql/queries.sql

COPY --from=build /go/src/github.com/max-samoylov/go-examples/fx-service/fx-service /srv/fx-service

ENTRYPOINT [ "/srv/wait.sh", "currencies-db:5432", "--", "/srv/fx-service" ]