# The TOP G Stack

This is a webapp that uses the TOP G stack (Tailwind, Ollama, Postgres, Golang).

TOP G allows you to create courses, add documents in those courses, and then
use LLMs to retrieve information from them.

## Quickstart

```console
docker build -t postgres-uuid .
docker run --rm -it -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 postgres-uuid
go run cmd/migrations/main.go init
go run cmd/migrations/main.go migrate
go run cmd/main.go # or install and run `air`
# open browser at `localhost:8080`
```
