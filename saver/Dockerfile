FROM golang:latest

ENV Saver_PATH = saver

RUN mkdir -p $GOPATH/src/${Saver_PATH}

WORKDIR $GOPATH/src/${Saver_PATH}

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . $GOPATH/src/${Saver_PATH}

RUN go build -o krishaSaver

CMD ["./krishaSaver"]