FROM golang:latest

RUN go get github.com/hectorj2f/search_web/web
RUN go install github.com/hectorj2f/search_web/web

CMD ["go", "run", "/go/src/github.com/hectorj2f/search_web/web/page.go"]

EXPOSE 8888
