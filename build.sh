#!/bin/sh

export VERSION=$GITHUB_REF_NAME
mkdir ./build

build() {
	EXTENSION=""
	if [ "$1" = "windows" ]; then
		EXTENSION=".exe"
	fi

	GOOS=$1 GOARCH=$2 go build -o ./build/ghcd$EXTENSION || exit 1
	cd ./build

	if [ "$1" != "windows" ]; then
		chmod u+x ./ghcd
	fi
	
	zip ./ghcd_${VERSION}_$1-$2.zip ./ghcd$EXTENSION || exit 1
	tar -czf ./ghcd_${VERSION}_$1-$2.tar.gz ./ghcd$EXTENSION || exit 1
	rm ./ghcd$EXTENSION
	cd ..
}

build linux amd64 || exit 1
build darwin amd64 || exit 1
build darwin arm64 || exit 1
build windows amd64 || exit 1
