#!/bin/bash
set -ex
cd `dirname $0`

RUN_NAME="apaas_ob_agent"

exec ./output/bin/$RUN_NAME
