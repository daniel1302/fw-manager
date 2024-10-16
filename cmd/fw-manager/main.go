package main

import (
	"flag"
	"fmt"
	"log"
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
	ipPOverride           string
}

var args fmArgs

func init() {
	flag.BoolVar(&args.withSSH, "with-ssh", true, "Decide if SSH port should be open. If false only hardcoded logic is applied")
	flag.BoolVar(&args.cleanAll, "clean-all", false, "Decide if all remaining rules should be removed")
	flag.BoolVar(&args.dryRun, "dry-run", false, "Decide if rules should be only printed to the output and not applied")
	flag.StringVar(&args.consulCatalogFilePath, "consul-catalog-file-path", "", "If not empty binary won't fetch catalog from consul API. Instead it will use given file")
	flag.StringVar(&args.networkCIDR, "network-cidr", "10.10.0.0/16", "The network CIDR for the wireguard")
	flag.StringVar(&args.ipPOverride, "ip-override", "", "If not empty program will assume local computer has assigned specific IP without checking it")
	flag.Parse()
}

func main() {
	normalizedFleetCatalog, err := normalizedCatalog(args.consulCatalogFilePath)
	if err != nil {
		log.Fatal("failed to get normalized fleet catalog", err)
	}

	thisComputerFleet, err := matchFleetServerToThisHost(args.ipPOverride, args.networkCIDR, normalizedFleetCatalog)
	if err != nil {
		log.Fatal("this computer does not belong to the managed network", err)
	}

	firewallRules := system.PrepareFirewallRules(*thisComputerFleet, &normalizedFleetCatalog)

	if args.dryRun {
		log.Println("Dry run only")

		for _, rule := range firewallRules {
			log.Printf("Open port %d for %s\n", rule.Port, rule.IP)
		}

		return
	}
}

func normalizedCatalog(consulCatalogFilePath string) (types.FleetCatalog, error) {
	var normalizedCatalog types.FleetCatalog
	if consulCatalogFilePath != "" {
		consulCatalog, err := consul.ReadLocalCatalog(consulCatalogFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read consul catalog from local file: %w", err)
		}

		normalizedCatalog, err = consul.NormalizeCatalog(consulCatalog)
		if err != nil {
			return nil, fmt.Errorf("failed to normalize local consul catalog: %w", err)
		}

		return normalizedCatalog, nil
	}

	consulApi, err := consul.NewConsulAPIClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul api client: %w", err)
	}

	consulDataCenters, err := consulApi.GetDataCenters()
	if err != nil {
		return nil, fmt.Errorf("failed to get data-centers from consul catalog: %w", err)
	}

	consulCatalog, err := consulApi.GetFleetCatalog(consulDataCenters)
	if err != nil {
		return nil, fmt.Errorf("failed to get fleet catalog from the consul api: %w", err)
	}
	normalizedCatalog, err = consul.NormalizeCatalog(consulCatalog)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize catalog ")
	}

	return normalizedCatalog, nil
}

func matchFleetServerToThisHost(ipOverride string, networkCIDR string, normalizedFleet types.FleetCatalog) (*types.FleetItem, error) {
	var (
		localIps []net.IP
		err      error
	)
	if ipOverride == "" {
		log.Printf("Getting IP addresses assigned to local interfaces")
		localIps, err = system.GetLocalIPs()
		if err != nil {
			return nil, fmt.Errorf("failed to get ip addresses assigned to local interfaces: %w", err)
		}
	} else {
		localIps = append(localIps, net.ParseIP(ipOverride))
	}

	_, wireguardCIDR, err := net.ParseCIDR(networkCIDR)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the network CIDR: %w", err)
	}

	var thisHostItem *types.FleetItem
	for _, localIP := range localIps {
		log.Printf("Checking IP: %s", localIP.String())

		if !wireguardCIDR.Contains(localIP) {
			log.Printf("... IP does not belong to wireguard cidr")
			continue
		}

		thisHostItem = normalizedFleet.FindItemByIP(localIP)
		if thisHostItem == nil {
			log.Println("... IP does not belong to the fleet, you want to manage")
		} else {
			log.Println("... IP belongs to the managed network")
			break
		}
	}

	if thisHostItem == nil {
		return nil, fmt.Errorf("this computer does not belong to the CIDR you specified")
	}

	return thisHostItem, nil
}
