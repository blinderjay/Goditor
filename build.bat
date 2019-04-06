#!/bin/sh
statik -src=app  -f -dest .  && \
go build .  && \
./Goditor