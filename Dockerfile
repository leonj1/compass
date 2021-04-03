FROM alpine:3.13.4
# using alpine instead of scratch since we need to mount a volume to the database

RUN apk update \
    && apk upgrade \
    && apk add bash

ADD compass /

ENTRYPOINT ["/compass"]

