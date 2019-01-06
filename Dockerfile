# build from latest go image
FROM golang:latest as build

WORKDIR /go/src/github.com/mchmarny/kgcs/
COPY . .

ENV GO111MODULE=on
RUN go mod download

WORKDIR /go/src/github.com/mchmarny/kgcs/cmd/service/

RUN CGO_ENABLED=0 go build -o /kgcs


# run image
FROM alpine as release
RUN apk add --no-cache ca-certificates

# app executable
COPY --from=build /kgcs /app/

# start server
WORKDIR /app
ENTRYPOINT ["./kgcs"]
