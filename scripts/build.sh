#!/usr/bin/env bash
# 
# Build datebook for release or local development
# This script will output binaries to the ./bin local directory
# 

# Get the parent directory of this script, which is the `datebook` repo
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

cd "$DIR"

if [[ $DATEBOOK_DEV = "1" ]]; then
  echo "building datebook in dev mode..."
  go build -o bin/datebook .
else
  echo "non-dev not supported yet..."
  exit 1
fi
