package version

import "fmt"

// Version information. Populated at build time.
var (
	Version   = "dev"
	Commit    = ""
	BuildDate = ""
)

// Print returns the formatted version string.
func Print() string {
	return fmt.Sprintf("%s (%s %s)", Version, Commit, BuildDate)
}
