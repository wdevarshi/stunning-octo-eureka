# Build Stage
FROM golang:1.24 AS build-stage

LABEL app="build-myapp"
LABEL REPO="https://github.com/ankurs/myapp"

ENV PROJPATH=/go/src/github.com/ankurs/myapp

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/ankurs/myapp
WORKDIR /go/src/github.com/ankurs/myapp

RUN make build-alpine

# Final Stage
FROM alpine:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/ankurs/myapp"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# add tz data
RUN apk add --no-cache tzdata

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/myapp/bin

WORKDIR /opt/myapp/bin

COPY --from=build-stage /go/src/github.com/ankurs/myapp/bin/myapp /opt/myapp/bin/
RUN chmod +x /opt/myapp/bin/myapp

# Create appuser
RUN adduser -D -g '' myapp
USER myapp

ENTRYPOINT ["/opt/myapp/bin/myapp"]
