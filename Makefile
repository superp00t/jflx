build:
	go build github.com/superp00t/jflx/cmd/jflx_server

all: build

install:
	cp jflx_server /usr/local/bin/jflx_server
