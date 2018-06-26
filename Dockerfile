FROM golang:1.10

ENV GOLANG_VERSION 1.10
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

ADD . /go/src/github.com/bialas1993/go-crawler

RUN go install github.com/bialas1993/go-crawler

WORKDIR $GOPATH

ENTRYPOINT /go/bin/go-crawler

EXPOSE 8080