#!/bin/bash

EXCLUDED="fileproto"
PACKAGES=$(go list ./... | grep -v "$EXCLUDED")

go test -v $PACKAGES -coverprofile=coverage.out > test_output.txt 2>&1
go tool cover -html=coverage.out 