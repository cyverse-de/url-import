FROM golang:1.16-alpine

RUN apk add --no-cache git
RUN go get -u github.com/jstemmer/go-junit-report

ARG git_commit=unknown
ARG descriptive_version=unknown

LABEL org.cyverse.git-ref="$git_commit"
LABEL org.cyverse.descriptive-version="$descriptive_version"

COPY . /go/src/github.com/cyverse-de/url-import
WORKDIR /go/src/github.com/cyverse-de/url-import
ENV CGO_ENABLED=0
RUN go install -v -ldflags="-X main.gitref=$git_commit" github.com/cyverse-de/url-import

LABEL org.label-schema.vcs-ref="$git_commit"
LABEL org.label-schema.vcs-url="https://github.com/cyverse-de/url-import"
LABEL org.label-schema.version="$descriptive_version"

ENTRYPOINT ["url-import"]
