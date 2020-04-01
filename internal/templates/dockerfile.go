// Package templates holds the template and project building functions
package templates

const dockerfileIdentifier = "dockerfile"

type dockerfile struct {
	*file
	Maintainer string
	StaticBin  string
}

// dockerfileConstructor returns the file content populated with the relevant values
func dockerfileConstructor(project *Project) (*file, error) {
	conf := project.Profile.Conf
	maintainer := conf.GetString("docker.maintainer")

	return newProjectFile(newDockerfile(maintainer, project.Name))
}

func newDockerfile(maintainer, staticbin string) *dockerfile {
	f, d, t := dockerfileValues()

	return &dockerfile{
		file:       newFile(dockerfileIdentifier, f, d, t),
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

func dockerfileValues() (f, d, t string) {
	const filename = "Dockerfile"

	const directory = "."

	const template = `# We could have used multi-stage,
# but compiling to a go static binary and inserting it is faster and smaller.

FROM gcr.io/distroless/static
LABEL maintainer="{{.Maintainer}}"
COPY {{.StaticBin}} {{.StaticBin}}
RUN echo "nonroot:x:65534:65534:nonroot:/:" > /etc/passwd
USER nonroot
ENTRYPOINT ["{{.StaticBin}}"]
`

	return filename, directory, template
}
