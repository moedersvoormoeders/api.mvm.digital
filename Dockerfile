FROM golang:1.14-alpine as build

COPY ./ /home/user/go/src/github.com/moedersvoormoeders/api.mvm.digital/

WORKDIR /home/user/go/src/github.com/moedersvoormoeders/api.mvm.digital/

ARG GOARCH
RUN GOARCH=${GOARCH} GOARM=7 go build ./cmd/mvmapi

FROM alpine:3.12

COPY --from=build /home/user/go/src/github.com/moedersvoormoeders/api.mvm.digital/mvmapi /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/mvmapi"]
CMD ["serve"]

