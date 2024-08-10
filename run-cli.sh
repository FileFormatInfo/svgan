#!/usr/bin/env bash
#
# run Jekyll locally
#

set -o errexit
set -o pipefail
set -o nounset

if [ -f ".env" ]; then
	export $(cat .env)
else
	echo "WARNING: no .env file found"
fi

if [ ! -d "tmp" ]; then
	mkdir tmp
fi

go build -o ./tmp/svgan ./cmd/svgan


tmp/svgan cmd/server/static/*.svg
