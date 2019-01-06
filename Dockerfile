# build from latest go image
FROM golang:latest as build

# copy
WORKDIR /go/src/github.com/mchmarny/gnotifs/
COPY . .

# modules
ENV GO111MODULE=on
RUN go mod download

# build
WORKDIR /go/src/github.com/mchmarny/gnotifs/cmd/service/
RUN CGO_ENABLED=0 go build -o /gnotifs


# run image
FROM alpine as release
RUN apk add --no-cache ca-certificates

# app executable
COPY --from=build /gnotifs /app/

# start server
WORKDIR /app
ENTRYPOINT ["./gnotifs"]
