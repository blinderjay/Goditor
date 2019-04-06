#!/bin/sh
statik -src=../app  -f -dest ../  && \
go install github.com/blinderjay/Goditor  && \
Goditor