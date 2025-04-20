#!/bin/bash

EXCLUDED="fileproto"
PACKAGES=$(go list ./... | grep -v "$EXCLUDED")

go test $PACKAGES -coverprofile=coverage.out -v > test_output.txt 2>&1
go tool cover -html=coverage.out 