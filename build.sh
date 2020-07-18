#!/bin/bash

# This fetches the k6-plugin-sql from GitHub and
# builds it in your $GOPATH/src/github.com/mostafa/k6-plugin-sql

go build -buildmode=plugin -ldflags="-s -w" -o sql.so github.com/mostafa/k6-plugin-sql
