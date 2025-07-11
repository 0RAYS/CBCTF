package utils

import (
	"fmt"
	"net"
)

func GetFirstIP(cidr string) (string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	ip := ipNet.IP.To4()
	if ip == nil {
		return "", fmt.Errorf("only IPv4 is supported in this example")
	}
	first := make(net.IP, len(ip))
	copy(first, ip)
	for i := len(first) - 1; i >= 0; i-- {
		first[i]++
		if first[i] != 0 {
			break
		}
	}
	if !ipNet.Contains(first) {
		return "", fmt.Errorf("no usable IPs in this CIDR")
	}
	return first.String(), nil
}
