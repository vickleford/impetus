FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN echo $GOPATH
COPY ./impetus /go/src/github.com/vickleford/impetus/impetus
RUN go get -d -v ./...
RUN go build -o impetuscli impetus/cmd/impetus/main.go
CMD mv impetuscli /artifacts
