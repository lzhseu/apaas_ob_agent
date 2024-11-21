#!/bin/bash
set -ex
cd `dirname $0`

RUN_NAME="apaas_ob_agent"

rm -rf output
mkdir -p output/bin output/conf
find conf/ -type f ! -name "*_local.*" | xargs -I{} cp {} output/conf/

go build -o output/bin/${RUN_NAME}
