#!/usr/bin/env zsh

dir="${0:A:h}"
tag="${tag:-$(git -C "${dir}" describe --tags)}"
commit="${commit:-$(git -C "${dir}" rev-parse HEAD)}"
date="${date:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}"

exec docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --build-arg VERSION="${tag}" \
    --build-arg COMMIT="${commit}" \
    --build-arg BUILD_DATE="${date}" \
    --tag leancodepublic.azurecr.io/kratos:"${tag}" \
    "${@}" "${dir}"
