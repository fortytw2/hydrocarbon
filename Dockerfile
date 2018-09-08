FROM golang:1.11-alpine as builder

RUN apk add yarn git bash

RUN go get -u github.com/lestrrat-go/bindata/...
RUN go get -u github.com/golang/dep/...

# Add our code
ADD ./ /go/src/github.com/fortytw2/hydrocarbon

# install node deps
WORKDIR /go/src/github.com/fortytw2/hydrocarbon/ui
RUN yarn install

# build
WORKDIR /go/src/github.com/fortytw2/hydrocarbon
RUN dep ensure
RUN go generate ./...
RUN CGO_ENABLED=0 go build -tags netgo -installsuffix netgo -o hydrocarbon github.com/fortytw2/hydrocarbon/cmd/hydrocarbon

# multistage
FROM alpine:latest

# https://stackoverflow.com/questions/33353532/does-alpine-linux-handle-certs-differently-than-busybox#33353762
RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=builder /go/src/github.com/fortytw2/hydrocarbon/hydrocarbon /usr/bin/hydrocarbon

# Run the image as a non-root user
RUN adduser -D hc
RUN chmod 0755 /usr/bin/hydrocarbon

USER hc

# Run the app. CMD is required to run on Heroku
CMD hydrocarbon