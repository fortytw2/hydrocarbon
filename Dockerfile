FROM golang:1.10-alpine as builder

RUN apk add yarn git bash

RUN go get -u github.com/lestrrat-go/bindata/...
RUN go get -u github.com/golang/dep/...

# Add our code
ADD ./ /go/src/github.com/fortytw2/hydrocarbon

# wipe git garbage if building locally
WORKDIR /go/src/github.com/fortytw2/hydrocarbon
RUN git clean -f -X

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

COPY --from=builder /go/src/github.com/fortytw2/hydrocarbon/hydrocarbon /usr/bin/hydrocarbon

# Run the image as a non-root user
RUN adduser -D hc
RUN chmod 0755 /usr/bin/hydrocarbon

USER hc

# Run the app. CMD is required to run on Heroku
CMD hydrocarbon