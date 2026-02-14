package utils

import (
	"fmt"
	"net/netip"
)

func GetLastIP(cidr string) (string, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return "", err
	}
	if !prefix.Addr().Is4() {
		return "", fmt.Errorf("only IPv4 is supported")
	}
	addr := prefix.Masked().Addr()
	raw := addr.As4()
	start := uint32(raw[0])<<24 | uint32(raw[1])<<16 | uint32(raw[2])<<8 | uint32(raw[3])
	num := uint32(1) << (32 - prefix.Bits())
	last := start + num - 2
	lastAddr := netip.AddrFrom4([4]byte{byte(last >> 24), byte(last >> 16), byte(last >> 8), byte(last)})
	if !prefix.Contains(lastAddr) {
		return "", fmt.Errorf("no usable IPs in this CIDR")
	}
	return lastAddr.String(), nil
}
