FROM golang:latest

ENV Ads_PATH = advertisement_scrapper

RUN mkdir -p $GOPATH/src/${Ads_PATH}

WORKDIR $GOPATH/src/${Ads_PATH}

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . $GOPATH/src/${Ads_PATH}

RUN go build -o krishaADS

CMD ["./krishaADS"]