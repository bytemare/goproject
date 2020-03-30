// Package templates holds the template and project building functions
package templates

type dockerfile struct {
	*file
	Maintainer string
	StaticBin  string
}

// dockerfileConstructor returns the file content populated with the relevant values
func dockerfileConstructor(project *Project) (*file, error) {
	conf := project.Profile.Conf
	maintainer := conf.GetString("author.name") + " " + conf.GetString("author.contact")

	return newProjectFile(newDockerfile(maintainer, project.Name))
}

func newDockerfile(maintainer, staticbin string) *dockerfile {
	return &dockerfile{
		file:       newFile(dockerfileIdentifier, dockerfileFilename, dockerTemplate),
		Maintainer: maintainer,
		StaticBin:  staticbin,
	}
}

func (d *dockerfile) getIdentifier() string {
	return d.identifier
}

func (d *dockerfile) getFilename() string {
	return d.filename
}

func (d *dockerfile) getTemplate() string {
	return d.template
}

const (
	dockerfileIdentifier = "dockerfile"
	dockerfileFilename   = "Dockerfile"
	dockerTemplate       = `# We could have used multi-stage,
# but compiling to a go static binary and inserting it is faster and smaller.

FROM gcr.io/distroless/static
LABEL maintainer="{{.Maintainer}}"
COPY {{.StaticBin}} {{.StaticBin}}
RUN echo "nonroot:x:65534:65534:nonroot:/:" > /etc/passwd
USER nonroot
ENTRYPOINT ["{{.StaticBin}}"]
`
)
