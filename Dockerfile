FROM golang:1.18

WORKDIR /build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go install -v ./...

RUN mv ./url-import /usr/local/bin/url-import 

ENTRYPOINT ["url-import"]
