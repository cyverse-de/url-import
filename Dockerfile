FROM golang:1.18

WORKDIR /build

ARG git_commit=unknown
ARG descriptive_version=unknown

LABEL org.cyverse.git-ref="$git_commit"
LABEL org.cyverse.descriptive-version="$descriptive_version"
LABEL org.label-schema.vcs-ref="$git_commit"
LABEL org.label-schema.vcs-url="https://github.com/cyverse-de/url-import"
LABEL org.label-schema.version="$descriptive_version"

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN go get -u github.com/jstemmer/go-junit-report

COPY . .

RUN go install -v ./...

ENTRYPOINT ["url-import"]
