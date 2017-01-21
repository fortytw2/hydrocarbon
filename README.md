Hydrocarbon [![Build Status](https://travis-ci.org/fortytw2/hydrocarbon.svg?branch=master)](https://travis-ci.org/fortytw2/hydrocarbon)
---------------------------------------------------------------------------------------------------------------------

> news, notifications, and updates. never miss out.

Plugin based news/feed reader with a focus on full text extraction and quality reading experiences on all platforms


# Dev Guide

You must have `lessc` and `go-bindata` installed and in your `$PATH`

Run `cmd/generate_cert` to generate a HTTPS cert in the `cmd/hydrocarbon` folder.

Then run a postgres on localhost:5432 ->

```
docker run -e "POSTGRES_PASSWORD=postgres" -d -p 5432:5432 postgres:9.6-alpine
```

then to rebuild and reload hydrocarbon run ->

```
cd $GOPATH/src/github.com/fortytw2/hydrocarbon/cmd/hydrocarbon && go generate github.com/fortytw2/hydrocarbon/... && go1.8beta2 build -tags dev && POSTGRES_DSN=postgres://postgres:postgres@localhost:5432?sslmode=disable ./hydrocarbon
```

# Screenshot

![Homepage Screenie](http://imgur.com/Ojktdiq.png)

LICENSE
-------

MIT, see LICENSE for full terms
