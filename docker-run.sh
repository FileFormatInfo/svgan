#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

APP_NAME="svgan-server"

echo "INFO: building docker image..."
docker build \
    --build-arg COMMIT=local@$(git rev-parse --short HEAD) \
    --build-arg LASTMOD=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
    --progress=plain \
    --tag "${APP_NAME}" \
    .

echo "INFO: running"
docker run \
	--publish 4000:4000 \
	--expose 4000 \
	--env PORT='4000' \
	"${APP_NAME}"

