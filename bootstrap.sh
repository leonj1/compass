#!/bin/bash

export HTTPPORT=${HTTPPORT:=80}

cd /app
/app/testify \
    -http-port=${HTTPPORT}

