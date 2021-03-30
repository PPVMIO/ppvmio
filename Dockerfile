FROM golang:1.15.7

RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/ppvmio
WORKDIR /go/src/ppvmio
ENV GOPATH /go
RUN dep ensure
CMD go run app.go
