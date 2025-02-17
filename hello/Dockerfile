FROM golang:1.23 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod tidy

COPY . /src

RUN CGO_ENABLED=0 go build -o /hello ./cmd/apiserver

FROM alpine:3.20

COPY --from=build /hello /hello

ENTRYPOINT [ "/hello" ]
