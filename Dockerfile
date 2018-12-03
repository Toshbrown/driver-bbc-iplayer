FROM amd64/alpine:3.8 as build
RUN echo http://nl.alpinelinux.org/alpine/edge/testing >> /etc/apk/repositories
RUN apk update && apk add build-base go git libzmq zeromq-dev alpine-sdk libsodium-dev

RUN addgroup -S databox && adduser -S -g databox databox
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/me-box/lib-go-databox
RUN go get -u golang.org/x/net/publicsuffix
COPY /src /src
COPY /static /static
RUN GGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo -ldflags '-s -w' -o driver /src/*.go

FROM amd64/alpine:3.8
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /driver /driver
COPY --from=build /static /static
RUN apk update && apk add libzmq && apk add ca-certificates
USER databox

LABEL databox.type="driver"
EXPOSE 8080

CMD ["./driver"]