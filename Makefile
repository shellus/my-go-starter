CURRENT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export GOPATH:=$(GOPATH)\;$(CURRENT)

default: window

test:
	echo $(GOPATH)

window:
	go install -ldflags "-H windowsgui" -tags 'release' win32api/createWindow

contributors:
	echo "Contributors to ngrok, both large and small:\n" > CONTRIBUTORS
	git log --raw | grep "^Author: " | sort | uniq | cut -d ' ' -f2- | sed 's/^/- /' | cut -d '<' -f1 >> CONTRIBUTORS
