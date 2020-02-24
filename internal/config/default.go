package config

const (
	DefaultProfileName    = "default.toml"
	DefaultProfileContent = `title = "default goproject profile"
	
	[author]
	name = "Bytemare"
	contact = "dev@bytema.re"
	
	[layout]
	directories = ["cmd", "internal"]
	files = [ "gitignore",
	        "doc",
	        "dockerfile",
	        "golangci",
	        "readme",
	        "sonar",
	        "travis"
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
	NewProfileContent = `title = ""
	
	[author]
	name = ""
	contact = ""
	
	[layout]
	directories = ["cmd", "internal"]
	files = [ "gitignore",
	        "doc",
	        "dockerfile",
	        "golangci",
	        "readme",
	        "sonar",
	        "travis"
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
