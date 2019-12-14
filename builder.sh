#!/bin/bash
docker run --rm -v "$PWD":/go/src/github.com/art-sitedesign/sitorama -w /go/src/github.com/art-sitedesign/sitorama golang:1.8 sh -c '
for GOOS in darwin linux; do
	for GOARCH in amd64; do
	  export GOOS GOARCH
	  go build -v -o ./bin/$GOOS-$GOARCH ./app
	done
done
'