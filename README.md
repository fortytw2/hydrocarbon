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

## license

mit
