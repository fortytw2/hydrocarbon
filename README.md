## hydrocarbon

> not just an rss reader.

## development

Run a copy of postgres somewhere with

```sh
docker run -p 5432:5432 postgres:alpine
```

```sh
go get -u github.com/fortytw2/hydrocarbon/...
cd $GOPATH/src/github.com/fortytw2/hydrocarbon/ui
yarn
cd $GOPATH/src/github.com/fortytw2/hydrocarbon/cmd/hydrocarbon
go generate $(go list ../../... | grep -v vendor) && go build -i . && POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable ./hydrocarbon -autoexplain
```

then open port :8080, enter an email, get the login token from hydrocarbon STDOUT
and proceed to develop.

## Configuring Image Server

Hydrocarbon has two modes for downloading and rehosting images, a local server
and a Google Cloud Storage backed server. To use the local server, configure nothing.

To configure google cloud storage, set `GCP_SERVICE_ACCOUNT`, `IMAGE_BUCKET_NAME`
and `IMAGE_DOMAIN`.

## license

mit
