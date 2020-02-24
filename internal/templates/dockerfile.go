package internal

type dockerfile struct {
	Maintainer string
	StaticBin  string
}

func newDockerfile(maintainer, staticbin string) *dockerfile {
	return &dockerfile{
		Maintainer: maintainer,
		StaticBin:  staticbin,
	}
}

func NewDockerfile(maintainer, staticbin string) (*ProjectFile, error) {
	return NewProjectFile("dockerfile", "Dockerfile", newDockerfile(maintainer, staticbin))
}

func (dockerfile) getTemplate() string {
	return `# We could have used multi-stage,
# but compiling to a go static binary and inserting it is faster and smaller.

FROM gcr.io/distroless/static
LABEL maintainer="{{.Maintainer}}"
COPY {{.StaticBin}} {{.StaticBin}}
RUN echo "nonroot:x:65534:65534:nonroot:/:" > /etc/passwd
USER nonroot
ENTRYPOINT ["{{.StaticBin}}"]
`
}
