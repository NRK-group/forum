FROM golang:1.19.8

RUN mkdir /forum

ADD . /forum

WORKDIR /forum

RUN go mod tidy

RUN cd /forum/cmd && go build  ./server.go

EXPOSE 8800

CMD ["/forum/cmd/server"]