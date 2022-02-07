#!/bin/bash
set -eux

BIN="golicense"
VERSION=$(git describe --tags)
CGO_ENABLED=0

CUR_DIR=$(cd "$(dirname "$0")" && pwd)
cd $CUR_DIR/..
TMP_DIR=dist/${BIN}

function package() {
	rm -rf ${TMP_DIR}
	mkdir -p ${TMP_DIR}
	cp LICENSE LICENSE-THIRD-PARTY ${TMP_DIR}
	go build -trimpath -ldflags "-w -s -X main.Version=${VERSION}" -o ${TMP_DIR}
	(
        cd "dist"
		if [ $GOOS = "windows" ]; then
			PKG=${BIN}-${VERSION}-${GOOS}-${GOARCH}.zip
			zip -r ${PKG} ${BIN}
		else
			PKG=${BIN}-${VERSION}-${GOOS}-${GOARCH}.tar.gz
			tar -zcf ${PKG} ${BIN}
		fi
		
        sha256sum ${PKG} > "${PKG}.sha256"
    )
}

GOOS=linux GOARCH=amd64 package
GOOS=linux GOARCH=arm64 package
GOOS=darwin GOARCH=amd64 package
GOOS=darwin GOARCH=arm64 package
GOOS=freebsd GOARCH=amd64 package
GOOS=windows GOARCH=amd64 package

(
	cd dist
	sha256sum -c --strict *.sha256
)