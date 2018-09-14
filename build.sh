#!/bin/bash

set -ex

# SET THE FOLLOWING VARIABLES
# docker hub username
USERNAME=fortisureapi
# image name
IMAGE=rest-api-go

version=`cat VERSION`

# Save the pwd before we run anything
PRE_PWD=`pwd`

# Determine the build script's actual directory, following symlinks
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
BUILD_DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"

# Derive the project name from the directory
PROJECT="$(basename $BUILD_DIR)"

# Setup the environment for the build
# GOPATH=$BUILD_DIR/src/vendor

# Build the project
cd $BUILD_DIR/src
GOOS=linux GOARCH=amd64 go build -o ../dist/$PROJECT *.go

EXIT_STATUS=$?

if [ $EXIT_STATUS == 0 ]; then
  echo "Build succeeded"

  echo "New binary in "$BUILD_DIR"/dist"

  cd $PRE_PWD

  docker build --no-cache -t $USERNAME/$IMAGE:latest .

  docker tag $USERNAME/$IMAGE:latest $USERNAME.azurecr.io/$IMAGE:$version
  docker tag $USERNAME/$IMAGE:latest $USERNAME.azurecr.io/$IMAGE:latest


  az acr login --name FortisureAPI

  docker push $USERNAME.azurecr.io/$IMAGE:$version
  docker push $USERNAME.azurecr.io/$IMAGE:latest

  EXIT_STATUS=$?

  if [ $EXIT_STATUS == 0 ]; then
    echo "Docker container created"
  else 
    echo "Docker container failed"
  fi

else
  echo "Build failed"
fi

# Change back to where we were
cd $PRE_PWD

exit $EXIT_STATUS
