package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
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

// GetRandomIP 从 CIDR 中随机选取一个 IPv4 地址
func GetRandomIP(cidr string) (string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", fmt.Errorf("invalid CIDR: %v", err)
	}

	ip = ip.To4()
	if ip == nil {
		return "", fmt.Errorf("only IPv4 is supported")
	}

	maskSize, bits := ipNet.Mask.Size()
	hostBits := bits - maskSize

	numHosts := new(big.Int).Lsh(big.NewInt(1), uint(hostBits))

	if numHosts.Cmp(big.NewInt(4)) < 0 {
		return "", fmt.Errorf("CIDR range too small to pick usable IPs")
	}

	// 可选的IP数量 = 总数 - 2（排除网络和广播地址）
	n, err := rand.Int(rand.Reader, new(big.Int).Sub(numHosts, big.NewInt(2)))
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %v", err)
	}
	n.Add(n, big.NewInt(1))

	networkIP := ip.Mask(ipNet.Mask)
	networkInt := big.NewInt(0).SetBytes(networkIP)
	resultInt := big.NewInt(0).Add(networkInt, n)
	result := resultInt.Bytes()
	for len(result) < 4 {
		result = append([]byte{0}, result...)
	}
	return net.IP(result).String(), nil
}

func ipToInt(ip net.IP) uint32 {
	ip4 := ip.To4()
	return uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3])
}

func intToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}
