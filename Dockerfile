ARG NODE_VERSION_IMAGE
FROM golang:1.22.2 as build
ARG VERSION_STRING="unknown"

COPY . /go

RUN cd /go && set -x; CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$VERSION_STRING" -o run

FROM ${NODE_VERSION_IMAGE}
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends openssl ca-certificates curl && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN chown -R node:node /app

COPY --from=build /go/run /app/run
RUN chmod +x /app/run

USER node

RUN sh -c 'whoami && /app/run true'

ENTRYPOINT ["/app/run"]
