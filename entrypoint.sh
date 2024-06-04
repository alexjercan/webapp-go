#!/bin/bash

./main db init && ./main db migrate;
# curl http://ollama:11434/api/pull -d '{"name": "llama3"}';
./main app run
