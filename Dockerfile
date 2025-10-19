# Build Stage
FROM golang:1.24 AS build-stage

LABEL app="transport-reliability-analytics"
LABEL description="LTA Transport Reliability Analytics System"

ENV PROJPATH=/go/src/app

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build-alpine

# Final Stage
FROM alpine:latest

ARG GIT_COMMIT
ARG VERSION
LABEL app="transport-reliability-analytics"
LABEL description="LTA Transport Reliability Analytics System"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Add timezone data and wget for health checks
RUN apk add --no-cache tzdata wget ca-certificates

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/transport-analytics/bin

WORKDIR /opt/transport-analytics/bin

COPY --from=build-stage /go/src/app/bin/myapp /opt/transport-analytics/bin/transport-analytics
RUN chmod +x /opt/transport-analytics/bin/transport-analytics

# Create non-root user for security
RUN adduser -D -g '' transport && \
    chown -R transport:transport /opt/transport-analytics

USER transport

EXPOSE 9090 9091

ENTRYPOINT ["/opt/transport-analytics/bin/transport-analytics"]
