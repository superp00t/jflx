FROM golang:1.25.3-alpine
WORKDIR /build

COPY ./ /build

RUN go build -v -o /build/jflx_server  github.com/superp00t/jflx/cmd/jflx_server
RUN go build -v -o /build/jflx_scan  github.com/superp00t/jflx/cmd/jflx_scan

FROM alpine
WORKDIR /jflx
COPY --from=0 /build/jflx_scan /bin/jflx_scan
COPY --from=0 /build/jflx_server /bin/jflx_server
CMD ["/bin/jflx_server"]
