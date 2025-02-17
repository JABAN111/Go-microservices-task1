FROM golang:1.23 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod tidy


COPY . /src


RUN CGO_ENABLED=0 go build -o /fileserver /src/cmd/apiserver/server.go

FROM alpine:3.20

COPY --from=build /fileserver /fileserver
RUN apk add curl

ENTRYPOINT [ "/fileserver" ]
