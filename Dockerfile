# syntax=docker/dockerfile:1
FROM --platform=${BUILDPLATFORM} docker.io/library/golang:alpine AS builder

ARG TARGETPLATFORM VERSION COMMIT BUILD_DATE

ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH="${TARGETPLATFORM#*/}"

WORKDIR /go/src/app

COPY go.mod go.sum ./
COPY internal/client-go/go.* internal/client-go/

RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
  go build \
  -ldflags="-s -w -X 'github.com/ory/kratos/driver/config.Version=${VERSION}' -X 'github.com/ory/kratos/driver/config.Commit=${COMMIT}' -X 'github.com/ory/kratos/driver/config.Date=${BUILD_DATE}'" \
  -o /kratos


FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder --chown=0:0 --chmod=755 /kratos /usr/bin/kratos

EXPOSE 4433 4434

ENTRYPOINT ["/usr/bin/kratos"]
CMD ["serve"]
