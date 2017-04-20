.PHONY: default server client deps fmt clean all release-all assets client-assets server-assets contributors
export GOPATH:=$(shell pwd)

default: client

client:
	go install -tags 'release' app/main/parse

contributors:
	echo "Contributors to ngrok, both large and small:\n" > CONTRIBUTORS
	git log --raw | grep "^Author: " | sort | uniq | cut -d ' ' -f2- | sed 's/^/- /' | cut -d '<' -f1 >> CONTRIBUTORS
