#!/bin/bash

# Get go dependencies for generate task first
go get github.com/go-openapi/runtime
go get golang.org/x/net/context
go get golang.org/x/net/context/ctxhttp
go get github.com/tylerb/graceful
go get github.com/jessevdk/go-flags

codon generate ${ARGS}