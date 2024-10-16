package system

import (
	"fmt"
	"net"
)

func GetLocalIPs() ([]net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to list network interfaces: %w", err)
	}

	result := []net.IP{}
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, fmt.Errorf("failed to list addresses assigned to the %s interface: %w", i.Name, err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			result = append(result, ip)
		}
	}

	return result, nil
}
