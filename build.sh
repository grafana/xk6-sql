#!/bin/bash

# This fetches the k6-plugin-<NAME> from GitHub and
# builds it in your $GOPATH/src/github.com/mostafa/k6-plugin-<NAME>

go build -buildmode=plugin -ldflags="-s -w" -o <name>.so github.com/<USERNAME>/k6-plugin-<NAME>
