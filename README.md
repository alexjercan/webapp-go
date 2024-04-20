# The TOP G Stack

This is a webapp that uses the TOP G stack (Tailwind, Ollama, Postgres, Golang).

TOP G allows you to create courses, add documents in those courses, and then
use LLMs to retrieve information from them.

## Quickstart

You will need to define the config.yaml and the secrets first. (They will be
kept in the repo for teaching purposes, just copy paste the example config and
fill in the github oauth stuff)

```console
docker-compose up
# open browser at `localhost:8080`
```
