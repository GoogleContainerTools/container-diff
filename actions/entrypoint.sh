#!/bin/sh -l

echo "$@"
sh -c "exec /go/bin/container-diff ${@}"
