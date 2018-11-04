FROM alpine as builder

RUN apk update && apk add go git gcc musl-dev
RUN mkdir ~/go
ENV GOPATH /go
WORKDIR $GOPATH

COPY * /go/src/

RUN go get -ldflags "-linkmode external -extldflags -static" -a instagram_clone
RUN go build -ldflags "-linkmode external -extldflags -static" -a instagram_clone

FROM scratch

COPY --from=builder /go/bin/instagram_clone /go/bin/instagram_clone
CMD ["/go/bin/instagram_clone"]
