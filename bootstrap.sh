#!/bin/bash

export HTTPPORT=${HTTPPORT:=80}

cd /app
/app/compass \
    -http-port=${HTTPPORT}

