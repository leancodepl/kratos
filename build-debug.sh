#!/usr/bin/env zsh

set -e

dir="${0:A:h}"
tag="${tag:-$(git -C "${dir}" describe --tags)}"
commit="${commit:-$(git -C "${dir}" rev-parse HEAD)}"
date="${date:-$(date -u +"%Y-%m-%d")}"
image="k3d-exampleapp-registry.local.lncd.pl:21345/kratos:${tag}-debug-$(date +%s)"

docker buildx build \
    --file Dockerfile.debug \
    --build-arg VERSION="${tag}" \
    --build-arg COMMIT="${commit}" \
    --build-arg BUILD_DATE="${date}" \
    --tag "${image}" \
    --load \
    "${dir}"

docker push "${image}"
docker rmi "${image}"

kubectl=(kubectl --context=k3d-exampleapp --namespace=kratos)

"${kubectl[@]}" patch deployment exampleapp-kratos --patch='{
  "spec": {
    "template": {
      "spec": {
        "containers": [
          {
            "name": "kratos",
            "image": "'"${image}"'",
            "resources": {
              "limits": {
                "cpu": "1",
                "memory": "1Gi"
              },
              "requests": {
                "cpu": "1",
                "memory": "1Gi"
              }
            }
          }
        ]
      }
    }
  }
}'

"${kubectl[@]}" rollout status deployments/exampleapp-kratos --watch=true

exec "${kubectl[@]}" port-forward deployments/exampleapp-kratos 2345
