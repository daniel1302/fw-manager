package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/daniel1302/fw-manager/consul"
	"github.com/daniel1302/fw-manager/system"
	"github.com/daniel1302/fw-manager/types"
)

type fmArgs struct {
	withSSH  bool
	cleanAll bool
	dryRun   bool

	consulCatalogFilePath string
	networkCIDR           string
}

var args fmArgs

func init() {
	flag.BoolVar(&args.withSSH, "with-ssh", true, "Decide if SSH port should be open. If false only hardcoded logic is applied")
	flag.BoolVar(&args.cleanAll, "clean-all", false, "Decide if all remaining rules should be removed")
	flag.BoolVar(&args.dryRun, "dry-run", false, "Decide if rules should be only printed to the output and not applied")
	flag.StringVar(&args.consulCatalogFilePath, "consul-catalog-file-path", "", "If not empty binary won't fetch catalog from consul API. Instead it will use given file")
	flag.StringVar(&args.networkCIDR, "network-cidr", "10.5.0.0/16", "The network CIDR for the wireguard")
	flag.Parse()
}

func main() {
	var normalizedCatalog types.FleetCatalog
	if args.consulCatalogFilePath != "" {
		consulCatalog, err := consul.ReadLocalCatalog(args.consulCatalogFilePath)
		if err != nil {
			panic(err)
		}

		normalizedCatalog, err = consul.NormalizeCatalog(consulCatalog)
		if err != nil {
			panic(err)
		}
	} else {
		consulApi, err := consul.NewConsulAPIClient(nil)
		if err != nil {
			panic(err)
		}

		consulDataCenters, err := consulApi.GetDataCenters()
		if err != nil {
			panic(err)
		}

		consulCatalog, err := consulApi.GetFleetCatalog(consulDataCenters)
		if err != nil {
			panic(err)
		}
		normalizedCatalog, err = consul.NormalizeCatalog(consulCatalog)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(normalizedCatalog)

	localIps, err := system.GetLocalIPs()
	if err != nil {
		panic(err)
	}

	fmt.Println(localIps)
	_, wireguardCIDR, err := net.ParseCIDR(args.networkCIDR)
	if err != nil {
		panic(err)
	}

	for _, localIP := range localIps {
		fmt.Printf("Checking IP: %s ... ", localIP.String())

		if !wireguardCIDR.Contains(localIP) {
			fmt.Printf("IP does not belong to wireguard cidr\n")
			continue
		}

		item := normalizedCatalog.FindItemByIP(localIP)
		if item == nil {
			fmt.Printf("Not found\n")
		} else {
			fmt.Printf("Found %s\n", item.Node)
		}
	}

}
