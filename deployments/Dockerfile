FROM golang:1.15

# in order to use cache.
RUN go get github.com/go-sql-driver/mysql

ADD .. /go/src/notes-web

WORKDIR /go/src/notes-web

RUN cp ../config/config.json /go/bin/ && \
    go install notes-web

ENTRYPOINT /go/bin/notes-web

EXPOSE 8888
