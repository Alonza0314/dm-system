#!/bin/bash

########################################################
# Script for building the docker image
#
# Usage:
#   ./build.sh <image tag>
#
# Description:
#   This script is used to build the docker image for dm-system.
#   The image tag is the name of the image to be built, default is latest.
########################################################

LATEST_TAG="latest"
IMAGE_NAME="alonza0314/dm-system"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

build_docker_image() {
    if ! docker build -f ${SCRIPT_DIR}/Dockerfile -t $IMAGE_NAME:$image_tag .; then
        echo "Failed to build the docker image"
        return 1
    fi
}

main() {
    local image_tag=${1:-$LATEST_TAG}

    if ! build_docker_image $image_tag; then
        return 1
    fi
}

main "$@"