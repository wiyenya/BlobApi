FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/BlobApi
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/BlobApi /go/src/BlobApi


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/BlobApi /usr/local/bin/BlobApi
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["BlobApi"]
