#FROM golang:1.10.3-alpine as gobuild
#WORKDIR /
#ENV GOPATH="/"
#RUN apk update && apk add pkgconfig build-base bash autoconf automake libtool gettext openrc git libzmq zeromq-dev mercurial
#COPY . . if you update the libs below build with --no-cache
FROM amd64/alpine:3.8 as build
RUN echo http://nl.alpinelinux.org/alpine/edge/testing >> /etc/apk/repositories
RUN apk update && apk add build-base go git libzmq zeromq-dev alpine-sdk libsodium-dev
RUN apk add 'go>=1.11-r0' --update-cache --repository http://nl.alpinelinux.org/alpine/edge/community

RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/me-box/lib-go-databox
RUN go get -u golang.org/x/net/publicsuffix
COPY /src /src
COPY /static /static
RUN addgroup -S databox && adduser -S -g databox databox
RUN GGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo -ldflags '-s -w' -o driver /src/*.go

#FROM alpine:3.8
#COPY --from=gobuild /etc/passwd /etc/passwd
FROM amd64/alpine:3.8
COPY --from=build /etc/passwd /etc/passwd
RUN apk update && apk add libzmq && apk add ca-certificates
USER databox
WORKDIR /
#COPY --from=gobuild /driver /driver
COPY --from=build /driver /driver
COPY --from=build /static /static
LABEL databox.type="driver"
EXPOSE 8080

CMD ["./driver"]
#CMD ["sleep","2147483647"]
