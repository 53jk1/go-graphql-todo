# build binary
FROM golang:1.19.4-alpine3.16 AS builder

RUN apk add --no-cache git

WORKDIR /go/src/github.com/53jk1/go-graphql-todo

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/go-graphql-todo ./cmd/server

FROM alpine:3.9

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /go/bin/go-graphql-todo /app/

EXPOSE 8080

CMD ["/app/go-graphql-todo"]