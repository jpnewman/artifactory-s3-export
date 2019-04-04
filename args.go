package main

import (
	"flag"
)

// Args - Arguments
type Args struct {
	updateS3Table *bool
	dryRun        *bool
}

func parseArgs() Args {
	r := Args{}
	r.updateS3Table = flag.Bool("updateS3Table", false, "Update S3 Table")
	r.dryRun = flag.Bool("dryrun", false, "Dry-Run")
	flag.Parse()

	return r
}
