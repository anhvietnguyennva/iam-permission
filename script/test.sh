#! /bin/bash

echo "Running test..."

go test -v -coverpkg=./... -coverprofile=profile.cov ./... > test_log.txt

sed -i '/^iam-permission\/cmd\/app\/main.go/ d' profile.cov
sed -i '/^iam-permission\/internal\/pkg\/migration\/migration.go/ d' profile.cov

go tool cover -func profile.cov
