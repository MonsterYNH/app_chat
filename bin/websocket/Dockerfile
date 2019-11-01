FROM golang:latest as builder
RUN mkdir /go/src/chat
COPY . /go/src/chat
RUN cd /go/src/chat/bin/websocket && go build main.go

FROM ubuntu:latest as WebSocket
COPY --from=builder /go/src/chat/bin/websocket /root
WORKDIR /root
EXPOSE 3456
CMD ./main