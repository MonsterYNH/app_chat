FROM golang:latest as builder
RUN mkdir /go/src/chat
COPY . /go/src/chat
RUN cd /go/src/chat/bin/static && go build main.go

FROM ubuntu:latest as Static
COPY --from=builder /go/src/chat/bin/static /root
WORKDIR /root
EXPOSE 6000
CMD ./main