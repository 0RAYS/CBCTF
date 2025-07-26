package traffic

import (
	"CBCTF/internal/i18n"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcap"
	"os"
)

type Connection struct {
	SrcIP   string
	DstIP   string
	SrcPort uint16
	DstPort uint16
	Type    string
}

func ReadPcap(path string) ([]Connection, bool, string) {
	if _, err := os.Stat(path); err != nil {
		return make([]Connection, 0), false, i18n.TrafficNotFound
	}
	handle, err := pcap.OpenOffline(path)
	if err != nil {
		return make([]Connection, 0), false, i18n.ReadPcapError
	}
	defer handle.Close()
	traffic := gopacket.NewPacketSource(handle, handle.LinkType())
	var connections []Connection
	for packet := range traffic.Packets() {
		var (
			srcIP   string
			dstIP   string
			srcPort uint16
			dstPort uint16
		)
		network := packet.NetworkLayer()
		if network == nil {
			continue
		}
		switch network.LayerType() {
		case layers.LayerTypeIPv4:
			if ipv4, ok := network.(*layers.IPv4); ok {
				srcIP = ipv4.SrcIP.String()
				dstIP = ipv4.DstIP.String()
			} else {
				continue
			}
		case layers.LayerTypeIPv6:
			if ipv6, ok := network.(*layers.IPv6); ok {
				srcIP = ipv6.SrcIP.String()
				dstIP = ipv6.DstIP.String()
			} else {
				continue
			}
		default:
			continue
		}
		transport := packet.TransportLayer()
		if transport == nil {
			continue
		}
		switch transport.LayerType() {
		case layers.LayerTypeTCP:
			if tcp, ok := transport.(*layers.TCP); ok {
				srcPort = uint16(tcp.SrcPort)
				dstPort = uint16(tcp.DstPort)
				connections = append(connections, Connection{
					SrcIP:   srcIP,
					DstIP:   dstIP,
					SrcPort: srcPort,
					DstPort: dstPort,
					Type:    layers.LayerTypeTCP.String(),
				})
			} else {
				continue
			}
		case layers.LayerTypeUDP:
			if udp, ok := transport.(*layers.UDP); ok {
				srcPort = uint16(udp.SrcPort)
				dstPort = uint16(udp.DstPort)
				connections = append(connections, Connection{
					SrcIP:   srcIP,
					DstIP:   dstIP,
					SrcPort: srcPort,
					DstPort: dstPort,
					Type:    layers.LayerTypeUDP.String(),
				})
			} else {
				continue
			}
		default:
			continue
		}
	}
	return connections, true, i18n.Success
}
