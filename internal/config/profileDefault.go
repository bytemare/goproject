// Package config groups the configuration and profile management mechanism
package config

const (
	defaultProfileName    = "default.toml"
	defaultProfileContent = `title = "default goproject profile"

	[author]
	name = "Bytemare"
	contact = "dev@bytema.re"

	[layout]
	directories = ["cmd", "internal", "internal/commands", "internal/config", "internal/version"]
	files = [ "gitignore",
	        "doc",
	        "dockerfile",
	        "golangci",
			"makefile",
			"pre-commit",
	        "readme",
	        "sonar",
	        "travis",
			"version"
		]

	[git]
	user = "bytemare"
	mail = "dev@bytema.re"
	URL  = "github.com/bytemare"

	[travis]
	profile = "https://travis-ci.com/bytemare"

	[sonar]
	org = "bytemare-github"

	[docker]
	maintainer = "Bytemare <dev@bytema.re>"
`
	newProfileContent = `title = ""

	[author]
	name = ""
	contact = ""

	[layout]
	directories = ["cmd", "internal", "internal/commands", "internal/config", "internal/version"]
	files = [ "gitignore",
	        "doc",
	        "dockerfile",
	        "golangci",
			"makefile",
			"pre-commit",
	        "readme",
	        "sonar",
	        "travis",
			"version"
		]

	[git]
	user = ""
	mail = ""
	URL  = ""

	[travis]
	profile = ""

	[sonar]
	org = ""

	[docker]
	maintainer = ""
`
)
