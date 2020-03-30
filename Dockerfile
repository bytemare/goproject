# 1. Build env
FROM golang:1.13.5-alpine AS build-envv
ADD goproject goproject

WORKDIR $GOPATH/src/github.com/bytemare/goproject/
COPY *.go ./
COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY .git ./
COPY Makefile ./
RUN make build-docker
#RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /bin/goproject ./cmd/goproject.go
    # && \
    #upx --ultra-brute -q /bin/goproject && \
    #upx -t /bin/goproject

# 2. Build image
FROM gcr.io/distroless/static
LABEL maintainer="Bytemare <dev@bytema.re>"
COPY --from=build-env /bin/goproject /bin/goproject
USER nonroot
ENTRYPOINT ["/bin/goproject"]
