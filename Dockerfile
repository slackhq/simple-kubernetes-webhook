# syntax=docker/dockerfile:1
# ---
FROM golang:1.27.2 AS build

ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR /work
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . /work

# Build admission-webhook
RUN --mount=type=cache,target=/root/.cache/go-build,sharing=private \
  --mount=type=cache,target=/go/pkg/mod \
  go build -o bin/admission-webhook .

# ---
FROM scratch AS run

COPY --from=build /work/bin/admission-webhook /usr/local/bin/

CMD ["admission-webhook"]
