FROM golang:1.19.5-alpine3.16
WORKDIR /couple-go
ADD . /couple-go
RUN go build -o couple-go
ENTRYPOINT ["./couple-go"]
