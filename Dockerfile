FROM golang:latest

RUN go get github.com/hectorj2f/search_web/web
RUN go install github.com/hectorj2f/search_web/web

ADD ./web/index.html /go/
RUN mkdir -p /go/static
COPY ./web/static /go/static

RUN cp /go/bin/web /go/search_web

CMD ["/go/search_web"]

EXPOSE 8888
