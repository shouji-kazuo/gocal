#!/bin/bash

set -e -u
package="$(basename $(pwd))"

if [ -e build ]; then
	rm -rf build
fi

mkdir build

gox -arch="amd64" -os="linux darwin windows" -output "build/${package}-{{.OS}}-{{.Arch}}"
