FROM golang:latest as builder
RUN mkdir /go/src/chat
COPY . /go/src/chat
RUN cd /go/src/chat/bin/api && go build main.go

FROM ubuntu:latest
COPY --from=builder /go/src/chat/bin/api /root
WORKDIR /root
EXPOSE 6000
CMD ./main