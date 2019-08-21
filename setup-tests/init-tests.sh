#!/bin/bash

set +x
# Grab the container definitions from the integration test
diffBase=$(cat ../tests/integration_test.go | grep diffBase | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
diffModified=$(cat ../tests/integration_test.go | grep diffModified | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
diffLayerBase=$(cat ../tests/integration_test.go | grep diffLayerBase | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
diffLayerModified=$(cat ../tests/integration_test.go | grep diffLayerModified | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
rpmBase=$(cat ../tests/integration_test.go | grep rpmBase | grep valentinrothberg | sed -e 's/.* = //g' | sed -e 's/"//g')
rpmModified=$(cat ../tests/integration_test.go | grep rpmModified | grep valentinrothberg | sed -e 's/.* = //g' | sed -e 's/"//g')
aptBase=$(cat ../tests/integration_test.go | grep aptBase | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
aptModified=$(cat ../tests/integration_test.go | grep aptModified | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
nodeBase=$(cat ../tests/integration_test.go | grep nodeBase | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
nodeModified=$(cat ../tests/integration_test.go | grep nodeModified | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
multiBase=$(cat ../tests/integration_test.go | grep -w multiBase | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
multiModified=$(cat ../tests/integration_test.go | grep -w multiModified | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
metadataBase=$(cat ../tests/integration_test.go | grep metadataBase | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
metadataModified=$(cat ../tests/integration_test.go | grep metadataModified | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')
pipModified=$(cat ../tests/integration_test.go | grep pipModified | grep gcr.io | sed -e 's/.* = //g' | sed -e 's/"//g')

# Echo the container definitions.
echo diffBase=$diffBase
echo diffModified=$diffModified
echo diffLayerBase=$diffLayerBase
echo diffLayerModified=$diffLayerModified
echo rpmBase=$rpmBase
echo rpmModified=$rpmModified
echo aptBase=$aptBase
echo aptModified=$aptModified
echo nodeBase=$nodeBase
echo nodeModified=$nodeModified
echo multiBase=$multiBase
echo multiModified=$multiModified
echo metadataBase=$metadataBase
echo metadataModified=$metadataModified
echo pipModified=$pipModified

# Now generate the containers
# XXX TO-DO for now we're only generated diffBase diffModified diffLayerBase and
# diffLayerModified, eventually we should generate all of them
#

echo 'docker build . -f Dockerfile.diffBase -t $diffBase:latest'
docker build . -f Dockerfile.diffBase -t $diffBase:latest
docker push $diffBase:latest
echo 'docker build . -f Dockerfile.diffLayerBase -t $diffLayerBase:latest'
docker build . -f Dockerfile.diffLayerBase -t $diffLayerBase:latest
docker push $diffLayerBase:latest
echo 'docker build . -f Dockerfile.diffLayerModified -t $diffLayerModified:latest'
docker build . -f Dockerfile.diffLayerModified -t $diffLayerModified:latest
docker push $diffLayerModified:latest
echo 'docker build . -f Dockerfile.diffModified -t $diffModified:latest'
docker build . -f Dockerfile.diffModified -t $diffModified:latest
docker push $diffModified:latest

#Now generate expected outputs.  Do NOT commit these without reviewing them for reasonableness
container-diff diff --no-cache -j --type=file $diffBase $diffModified > ../tests/file_diff_expected.json
container-diff diff --no-cache -j --type=layer $diffLayerBase $diffLayerModified > ../tests/file_layer_diff_expected.json
container-diff diff --no-cache -j --type=size $diffLayerBase $diffLayerModified > ../tests/size_diff_expected.json
container-diff diff --no-cache -j --type=sizelayer $diffLayerBase $diffLayerModified > ../tests/size_layer_diff_expected.json
container-diff diff --no-cache -j --type=apt $aptBase $aptModified > ../tests/apt_diff_expected.json
container-diff diff --no-cache -j --type=node $nodeBase $nodeModified > ../tests/node_diff_order_expected.json
container-diff diff --no-cache -j --type=node --type=pip --type=apt $multiBase $multiModified > ../tests/multi_diff_expected.json
container-diff diff --no-cache -j --type=history $diffBase $diffModified > ../tests/hist_diff_expected.json
container-diff diff --no-cache -j --type=metadata $metadataBase $metadataModified > ../tests/metadata_diff_expected.json
container-diff diff --no-cache -j --type=apt -o $aptBase $aptModified > ../tests/apt_sorted_diff_expected.json
container-diff analyze --no-cache -j --type=apt $aptModified > ../tests/apt_analysis_expected.json
container-diff analyze --no-cache -j --type=file -o $diffModified > ../tests/file_sorted_analysis_expected.json
container-diff analyze --no-cache -j --type=layer $diffLayerBase > ../tests/file_layer_analysis_expected.json
container-diff analyze --no-cache -j --type=size $diffBase > ../tests/size_analysis_expected.json
container-diff analyze --no-cache -j --type=sizelayer $diffLayerBase > ../tests/size_layer_analysis_expected.json
container-diff analyze --no-cache -j --type=pip $pipModified > ../tests/pip_analysis_expected.json
container-diff analyze --no-cache -j --type=node $nodeModified > ../tests/node_analysis_expected.json

