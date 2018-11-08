FROM golang:1.10.3-alpine as gobuild
WORKDIR /
ENV GOPATH="/"
RUN apk update && apk add pkgconfig build-base bash autoconf automake libtool gettext openrc git libzmq zeromq-dev mercurial
#COPY . . if you update the libs below build with --no-cache
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/me-box/lib-go-databox
RUN go get -u golang.org/x/net/publicsuffix
COPY . .
RUN addgroup -S databox && adduser -S -g databox databox
RUN GGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo -ldflags '-s -w' -o driver /src/*.go

FROM alpine:3.8
COPY --from=gobuild /etc/passwd /etc/passwd
RUN apk update && apk add libzmq
USER databox
WORKDIR /
COPY --from=gobuild /driver .
LABEL databox.type="driver"
EXPOSE 8080

CMD ["./driver"]
#CMD ["sleep","2147483647"]
