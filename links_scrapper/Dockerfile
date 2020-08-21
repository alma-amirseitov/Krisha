FROM golang:latest

ENV Links_PATH = links_scrapper

RUN mkdir -p $GOPATH/src/${Links_PATH}

WORKDIR $GOPATH/src/${Links_PATH}

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . $GOPATH/src/${Links_PATH}

RUN go build -o krishaURL

CMD ["./krishaURL"]