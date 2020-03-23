FROM alpine:3.8

EXPOSE 9091

RUN apk add --no-cache ca-certificates libc6-compat

COPY ./bin/server /server

ENTRYPOINT [ "/server" ]
