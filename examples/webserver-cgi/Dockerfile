FROM ubuntu:16.04

RUN apt update && apt install -y --no-install-recommends fcgiwrap && rm -rf /var/cache/apk/*

STOPSIGNAL SIGTERM
ENTRYPOINT [ "fcgiwrap", "-s", "tcp:0.0.0.0:9000", "-f" ]
