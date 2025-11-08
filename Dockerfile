FROM golang:1.25.3-alpine@sha256:aee43c3ccbf24fdffb7295693b6e33b21e01baec1b2a55acc351fde345e9ec34

WORKDIR /build

COPY ./ /build

RUN go build -v -o /build/jflx_server  github.com/superp00t/jflx/cmd/jflx_server

FROM alpine
COPY --from=0 /build/jflx_server /bin/jflx_server
CMD ["/bin/jflx_server"]
