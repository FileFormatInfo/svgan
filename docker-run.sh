#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

docker build -t svgan-server .

echo "INFO: running"
docker run \
	--publish 4000:4000 \
	--expose 4000 \
	--env PORT='4000' \
	--env LASTMOD=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
	svgan-server

	#--env COMMIT=$(git rev-parse --short HEAD) \
