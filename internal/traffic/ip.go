package traffic

import (
	"net/netip"
	"strings"
)

// Special-purpose ranges that IsGlobalUnicast alone does not exclude.
var specialTrafficPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"), netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("192.0.0.0/24"), netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.88.99.0/24"), netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"), netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("240.0.0.0/4"), netip.MustParsePrefix("::/96"),
	netip.MustParsePrefix("64:ff9b::/96"), netip.MustParsePrefix("64:ff9b:1::/48"),
	netip.MustParsePrefix("100::/64"), netip.MustParsePrefix("2001::/23"),
	netip.MustParsePrefix("2001:db8::/32"), netip.MustParsePrefix("2002::/16"),
	netip.MustParsePrefix("3fff::/20"), netip.MustParsePrefix("5f00::/16"),
	netip.MustParsePrefix("fec0::/10"),
}

func NormalizeTrafficIP(value string) string {
	ip, err := netip.ParseAddr(strings.TrimSpace(value))
	if err != nil {
		return ""
	}
	return ip.WithZone("").Unmap().String()
}

func IsPublicTrafficIP(value string) bool {
	ip, err := netip.ParseAddr(NormalizeTrafficIP(value))
	if err != nil || !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() {
		return false
	}
	for _, prefix := range specialTrafficPrefixes {
		if prefix.Contains(ip) {
			return false
		}
	}
	return true
}
