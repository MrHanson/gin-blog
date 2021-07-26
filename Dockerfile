FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/https://github.com/MrHanson/gin-blog
COPY . $GOPATH/src/https://github.com/MrHanson/gin-blog
RUN go build .

EXPOSE 8000
ENTRYPOINT ["./gin-blog"]
