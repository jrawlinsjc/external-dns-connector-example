ARG PUBLIC_ECR_REPOSITORY_PATH=public.ecr.aws/docker/library
ARG GO_BASE_IMAGE=golang
ARG GO_BASE_IMAGE_VERSION=1.23.2

FROM ${PUBLIC_ECR_REPOSITORY_PATH}/${GO_BASE_IMAGE}:${GO_BASE_IMAGE_VERSION} AS base
WORKDIR /external-dns-connector-example

ENV GOFLAGS=-mod=readonly

COPY . .
RUN go build -o external-dns-connector-example .

FROM busybox
WORKDIR /

COPY --from=base /external-dns-connector-example/external-dns-connector-example .

CMD ["/external-dns-connector-example"]
