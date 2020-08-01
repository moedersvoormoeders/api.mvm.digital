FROM golang:1.14-alpine as build

COPY ./ /home/user/go/src/github.com/moedersvoormoeders/api.mvm.digital/

WORKDIR /home/user/go/src/github.com/moedersvoormoeders/api.mvm.digital/

ARG GOARCH
RUN GOARCH=${GOARCH} GOARM=7 go build ./cmd/mvmapi

ARG ARCH
FROM $ARCH/alpine:3.12

ARG QEMU_BIN
COPY $QEMU_BIN /usr/bin

RUN apk add --no-cache ca-certificates

COPY --from=build /home/user/go/src/github.com/moedersvoormoeders/api.mvm.digital/mvmapi /usr/local/bin/

RUN mkdir /opt/mvm-api
WORKDIR /opt/mvm-api

ENTRYPOINT ["/usr/local/bin/mvmapi"]
CMD ["serve"]

