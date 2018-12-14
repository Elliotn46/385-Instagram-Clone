FROM alpine as builder

RUN apk update && apk add go git gcc musl-dev graphicsmagick
RUN mkdir ~/go
ENV GOPATH /go
WORKDIR $GOPATH

COPY * /go/src/

RUN go get -ldflags "-linkmode external -extldflags -static" -a rabbit
RUN go build -ldflags "-linkmode external -extldflags -static" -a rabbit

FROM alpine
RUN apk update && apk add graphicsmagick
COPY --from=builder /go/bin/rabbit /go/bin/rabbit
CMD ["/go/bin/rabbit"]
