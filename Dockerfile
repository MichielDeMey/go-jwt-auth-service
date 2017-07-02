FROM golang:1-alpine

WORKDIR /go/src/michiel.be/jwt-auth

COPY . .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]
