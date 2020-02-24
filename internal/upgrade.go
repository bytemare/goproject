package internal

import "fmt"

var (
	version string
	commit  string
)

func PrintableVersion() string {
	return fmt.Sprintf("Version %s:%s", version, commit)
}

func Upgrade() {
	fmt.Printf("%s - Upgrade not implemented yet.\n", PrintableVersion())
}
