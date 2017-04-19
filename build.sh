#!/bin/bash -x

go get -u github.com/alecthomas/gometalinter \
&& go install henryleong.com/... \
&& GOOS=windows go build henryleong.com/theserver/... \
&& GOOS=windows go build henryleong.com/singlesource/... \

