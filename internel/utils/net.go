package utils

import (
	"fmt"
	"net"
)

var (
	ipBlockL = make([][]string, 0)
	ipBlockN = make(map[string]uint)
)

func GetIPBlock(n uint, cidr string, blockSize int) ([]string, error) {
	var err error
	if len(ipBlockL) == 0 {
		ipBlockL, err = splitCIDR(cidr, blockSize)
		if err != nil {
			return nil, err
		}
	}
	block := ipBlockL[n%uint(len(ipBlockL))]
	retry := 0
	for {
		retry++
		if retry > len(ipBlockL) {
			return make([]string, 0), fmt.Errorf("no available IP block")
		}
		key := fmt.Sprintf("%s-%d", block[0], len(block))
		if _, ok := ipBlockN[key]; ok {
			n++
			block = ipBlockL[n%uint(len(ipBlockL))]
			continue
		}
		ipBlockN[key] = n
		break
	}
	return block, nil
}

func RemoveIPBlock(block string) {
	if _, ok := ipBlockN[block]; ok {
		delete(ipBlockN, block)
	}
}

func splitCIDR(cidr string, blockSize int) ([][]string, error) {
	blocks := make([][]string, 0)
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return blocks, err
	}
	maskSize, _ := ipNet.Mask.Size()
	if blockSize < maskSize || blockSize > 32 {
		return blocks, fmt.Errorf("block size must be between %d and 32", maskSize)
	}
	blockCount := 1 << (blockSize - maskSize)
	blockIPCount := 1 << (32 - blockSize)
	startIP := ipNet.IP.To4()
	if startIP == nil {
		return blocks, fmt.Errorf("not IPv6")
	}
	for i := 0; i < blockCount; i++ {
		blockStart := make(net.IP, len(startIP))
		copy(blockStart, startIP)
		offset := i * blockIPCount
		for j := 0; j < offset; j++ {
			incrementIP(blockStart)
		}
		currentIP := make(net.IP, len(blockStart))
		copy(currentIP, blockStart)
		ipL := []string{blockStart.String()}
		for j := 0; j < blockIPCount-1; j++ {
			incrementIP(currentIP)
			ipL = append(ipL, currentIP.String())
		}
		blocks = append(blocks, ipL)
	}
	return blocks, nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
