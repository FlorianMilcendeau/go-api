FROM golang:alpine

WORKDIR /usr/src

ADD . .

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build /usr/src/main.go" --command="./main"