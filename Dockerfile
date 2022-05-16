FROM golang:1.17

WORKDIR /forum

COPY . .

# RUN go mod init forum/forum

RUN go mod tidy

RUN cd /forum/cmd && go build  ./server.go

EXPOSE 8800

CMD ["/forum/cmd/server"]