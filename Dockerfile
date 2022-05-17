FROM golang:1.17

RUN mkdir /forum

ADD . /forum

WORKDIR /forum

RUN go mod tidy

RUN cd /forum/cmd && go build  ./server.go

EXPOSE 8800

CMD ["/forum/cmd/server"]