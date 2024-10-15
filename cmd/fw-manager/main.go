package main

import (
	"flag"

	"github.com/daniel1302/fw-manager/consul"
)

type fmArgs struct {
	withSSH       bool
	cleanAll      bool
	consulCatalog bool
	dryRun        bool
}

var args fmArgs

func init() {
	flag.BoolVar(&args.withSSH, "with-ssh", true, "Decide if SSH port should be open. If false only hardcoded logic is applied")
	flag.BoolVar(&args.cleanAll, "clean-all", false, "Decide if all remaining rules should be removed")
	flag.BoolVar(&args.dryRun, "dry-run", false, "Decide if rules should be only printed to the output and not applied")

	flag.Parse()
}

func main() {
	consulApi, err := consul.NewConsulAPIClient(nil)
	if err != nil {
		panic(err)
	}

	consulDataCenters, err := consulApi.GetDataCenters()
	if err != nil {
		panic(err)
	}

	consulApi.GetFleetCatalog(consulDataCenters)
}
