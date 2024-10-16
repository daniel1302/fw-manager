package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/daniel1302/fw-manager/system"
	consulapi "github.com/hashicorp/consul/api"
)

var httpSvcPost int
var svcName string
var svcEnv string
var svcStage string
var svcAddr string

func init() {
	flag.IntVar(&httpSvcPost, "port", 8011, "The http port for the service")
	flag.StringVar(&svcName, "name", "wireguard", "The service name reported to consul")
	flag.StringVar(&svcEnv, "env", "metrics", "The service env reported to consul in the metadata section")
	flag.StringVar(&svcStage, "stage", "test", "The service stage reported to consul in the metadata section")
	flag.StringVar(&svcAddr, "address", "localhost", "The service address")

	flag.Parse()
}

func main() {
	//nolint:errcheck
	register()

	log.Fatal(runHTTPServer(httpSvcPost))
}

func runHTTPServer(port int) error {
	mux := http.NewServeMux()
	mux.Handle("/", &SvcHandler{})

	log.Printf("HTTP server is listening on :%d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

type SvcHandler struct {
	counter int
}

func (svc *SvcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Nothing interesting is here :) #%d", svc.counter)
	svc.counter++
}

func findLocalIP(cidr string) (*net.IP, error) {
	localIPs, err := system.GetLocalIPs()
	if err != nil {
		return nil, fmt.Errorf("failed to get local ip addresses: %w", err)
	}

	_, localCIDR, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cidr: %w", err)
	}

	for idx, localIp := range localIPs {
		if localCIDR.Contains(localIp) {
			return &localIPs[idx], nil
		}
	}

	// todo: We really do not want this service to fail. In future we should create custom error.
	return nil, nil
}

func register() error {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create consul client: %w", err)
	}

	localIp, err := findLocalIP("10.10.0.0/16")
	if err != nil {
		localIp = &[]net.IP{net.ParseIP("10.10.0.10")}[0]
	}

	registeration := &consulapi.AgentServiceRegistration{
		ID:      svcName,
		Name:    svcName,
		Port:    httpSvcPost,
		Address: localIp.String(),
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%v/check", svcAddr, httpSvcPost),
			Interval: "10s",
			Timeout:  "30s",
		},
		Meta: map[string]string{
			"env":   svcEnv,
			"stage": svcStage,
		},
		Tags: []string{
			fmt.Sprintf("%s.%s", svcEnv, svcStage),
			"wireguard",
			"vpn",
		},
	}

	if err = consul.Agent().ServiceRegister(registeration); err != nil {
		return fmt.Errorf("failed to register service(%s:%v): %w", svcName, httpSvcPost, err)
	} else {
		log.Printf("successfully register service: %s:%v", svcName, httpSvcPost)
	}

	return nil
}
