#!/bin/bash -x

go get -u github.com/alecthomas/gometalinter \
&& go get -u github.com/gorilla/mux \
&& go install henryleong.com/...

