/**
 * This is a standard Cobra CLI app.
 * We don't really use subcommands for now, only the Root one.
 *
 * A good entrypoint is therefore cmd/root.go
 */

package main

import "github.com/MieuxVoter/majority-judgment-cli/cmd"

func main() {
	cmd.Execute()
}
