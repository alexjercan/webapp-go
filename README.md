# WebApp in Go

Simple web app in go + postgresql + htmx to learn go.

## Quickstart

```console
docker build -t postgres-uuid .
docker run --rm -it -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 postgres-uuid
go run cmd/main.go # or install and run `air`
# open browser at `localhost:8080`
```
