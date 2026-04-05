package ratelimit

import "net/netip"

type Allowlist struct {
	addrs    map[netip.Addr]struct{}
	prefixes []netip.Prefix
}

func NewAllowlist(entries []string) Allowlist {
	allowlist := Allowlist{
		addrs: make(map[netip.Addr]struct{}, len(entries)),
	}
	for _, entry := range entries {
		if prefix, err := netip.ParsePrefix(entry); err == nil {
			allowlist.prefixes = append(allowlist.prefixes, prefix)
			continue
		}
		if addr, err := netip.ParseAddr(entry); err == nil {
			allowlist.addrs[addr.Unmap()] = struct{}{}
		}
	}
	return allowlist
}

func (a Allowlist) Contains(rawIP string) bool {
	addr, err := netip.ParseAddr(rawIP)
	if err != nil {
		return false
	}
	addr = addr.Unmap()
	if _, ok := a.addrs[addr]; ok {
		return true
	}
	for _, prefix := range a.prefixes {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}
