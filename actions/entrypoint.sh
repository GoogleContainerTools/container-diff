#!/bin/bash

echo "$@"
/go/bin/container-diff ${@}
