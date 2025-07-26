package utils

import (
	"fmt"
	"net"
)

func GetGatewayIP(cidr string) (string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	ip := ipNet.IP.To4()
	if ip == nil {
		return "", fmt.Errorf("only IPv4 is supported")
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

func GetLastIP(cidr string) (string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return "", fmt.Errorf("only IPv4 is supported")
	}
	start := ipToInt(ip.Mask(ipNet.Mask))
	ones, bits := ipNet.Mask.Size()
	num := uint32(1) << uint32(bits-ones)
	lastIP := intToIP(start + num - 2)
	if !ipNet.Contains(lastIP) {
		return "", fmt.Errorf("no usable IPs in this CIDR")
	}
	return lastIP.String(), nil
}

func ipToInt(ip net.IP) uint32 {
	ip4 := ip.To4()
	return uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3])
}

func intToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}
