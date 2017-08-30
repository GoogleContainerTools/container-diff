#!/bin/bash
while IFS=$' \n\r' read -r flag differ image1 image2 file; do
  go run main.go diff $image1 $image2 $flag -j > $file
  if [[ $? -ne 0 ]]; then
    echo "ERROR container-diff diff $differ differ failed"
    exit 1
  fi
done < tests/test_differ_runs.txt

while IFS=$' \n\r' read -r flag analyzer image file; do
  go run main.go analyze $image $flag -j > $file
  if [[ $? -ne 0 ]]; then
    echo "ERROR: container-diff analyze $analyzer analyzer failed"
    exit 1
  fi
done < tests/test_analyzer_runs.txt

success=0
while IFS=$' \n\r' read -r type analyzer actual expected; do
  diff=$(jq --argfile a "$actual" --argfile b "$expected" -n 'def walk(f): . as $in | if type == "object" then reduce keys[] as $key ( {}; . + { ($key):  ($in[$key] | walk(f)) } ) | f elif type == "array" then map( walk(f) ) | f else f end; ($a | walk(if type == "array" then sort else . end)) as $a | ($b | walk(if type == "array" then sort else . end)) as $b | $a == $b')
  if ! "$diff" ; then
    echo "ERROR: container-diff analyze $analyzer $type: output is not as expected"
    success=1
  fi
done < tests/test_run_comparisons.txt
if [[ "$success" -ne 0 ]]; then
  exit 1
fi

go test `go list ./... | grep -v vendor`
