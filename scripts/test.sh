#!/bin/bash

EXCLUDED="fileproto"
PACKAGES=$(go list ./... | grep -v "$EXCLUDED")

cd web 
npm i
npm run build 
cd ..

tree

go test -v $PACKAGES -coverprofile=coverage.out > test_output.txt 2>&1
go tool cover -html=coverage.out 