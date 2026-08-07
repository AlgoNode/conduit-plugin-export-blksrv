package main

import (
	"fmt"
	"os"

	// Imports for built-in plugins.
	// Exporters are imported individually instead of via exporters/all to
	// avoid the postgresql exporter's unfixed pgx/v4 vulnerability chain.
	_ "github.com/algorand/conduit/conduit/plugins/exporters/filewriter"
	_ "github.com/algorand/conduit/conduit/plugins/exporters/noop"
	_ "github.com/algorand/conduit/conduit/plugins/importers/all"
	_ "github.com/algorand/conduit/conduit/plugins/processors/all"

	_ "github.com/algorand/conduit-plugin-template/plugin/exporter"

	"github.com/algorand/conduit/pkg/cli"
)

func main() {
	conduitCmd := cli.MakeConduitCmdWithUtilities()
	if err := conduitCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
