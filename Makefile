CURRENT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export GOPATH:=$(GOPATH)\;$(CURRENT)

default: bindata

test:
	echo $(GOPATH)

bindata:
	go-bindata -o src/bindata/bindata.go -prefix src/bindata/ src/bindata/asset/...
	go install -tags 'debug' bindata

window:
	go install -ldflags "-H windowsgui" -tags 'release' win32api/createWindow

contributors:
	echo "Contributors to ngrok, both large and small:\n" > CONTRIBUTORS
	git log --raw | grep "^Author: " | sort | uniq | cut -d ' ' -f2- | sed 's/^/- /' | cut -d '<' -f1 >> CONTRIBUTORS
