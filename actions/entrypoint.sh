#!/bin/bash

command="${INPUT_COMMAND} ${INPUT_ARGS}"
echo "container-diff ${command}"
/usr/local/bin/container-diff ${command}
