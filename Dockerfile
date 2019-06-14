FROM golang
COPY . /webserver
WORKDIR /webserver
CMD go run app.go
