package utils

import (
	"CBCTF/internal/config"
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcap"
	pp "github.com/pires/go-proxyproto"
)

type Connection struct {
	TimeShift time.Duration
	Time      time.Time
	SrcIP     string
	DstIP     string
	Type      string
	Subtype   string
	Size      int
}

func ReadPcapFile(path string) ([]Connection, error) {
	file, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if file.IsDir() {
		return nil, fmt.Errorf("%s is a directory", path)
	}
	handle, err := pcap.OpenOffline(path)
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	traffic := gopacket.NewPacketSource(handle, handle.LinkType())
	var connections []Connection
	var firstPacketTime time.Time
	for packet := range traffic.Packets() {
		connection := Connection{Size: packet.Metadata().CaptureLength, Time: packet.Metadata().Timestamp}
		if firstPacketTime.IsZero() {
			firstPacketTime = packet.Metadata().Timestamp
		}
		connection.TimeShift = packet.Metadata().Timestamp.Sub(firstPacketTime)
		network := packet.NetworkLayer()
		if network == nil {
			continue
		}
		switch network.LayerType() {
		case layers.LayerTypeIPv4:
			if ipv4, ok := network.(*layers.IPv4); ok {
				connection.SrcIP = ipv4.SrcIP.String()
				connection.DstIP = ipv4.DstIP.String()
			} else {
				continue
			}
		case layers.LayerTypeIPv6:
			if ipv6, ok := network.(*layers.IPv6); ok {
				connection.SrcIP = ipv6.SrcIP.String()
				connection.DstIP = ipv6.DstIP.String()
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
			connection.Type = layers.LayerTypeTCP.String()
		case layers.LayerTypeUDP:
			connection.Type = layers.LayerTypeUDP.String()
		default:
			continue
		}
		application := packet.ApplicationLayer()
		if application == nil {
			continue
		}
		connection.Subtype = application.LayerType().String()
		if config.Env.K8S.Frpc.On {
			if header, err := pp.Read(bufio.NewReader(bytes.NewReader(transport.LayerPayload()))); err == nil {
				srcIP, _, err := net.SplitHostPort(header.SourceAddr.String())
				if err != nil {
					continue
				}
				dstIP, _, err := net.SplitHostPort(header.DestinationAddr.String())
				if err != nil {
					continue
				}
				connection.SrcIP = srcIP
				connection.DstIP = dstIP
				connection.Subtype = "Proxy"
			}
		}
		connections = append(connections, connection)
	}
	return connections, nil
}

func ReadPcapDir(path string) ([]Connection, error) {
	d, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !d.IsDir() {
		return nil, fmt.Errorf("%s is a file", path)
	}
	connections := make([]Connection, 0)
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range dir {
		if file.IsDir() || (!strings.HasSuffix(file.Name(), ".pcap") && !strings.HasSuffix(file.Name(), ".pcapng")) {
			continue
		}
		packet, err := ReadPcapFile(fmt.Sprintf("%s/%s", path, file.Name()))
		if err != nil {
			return nil, err
		}
		connections = append(connections, packet...)
	}
	if len(connections) < 1 {
		return nil, nil
	}
	slices.SortStableFunc(connections, func(c1 Connection, c2 Connection) int { return c1.Time.Compare(c2.Time) })
	firstPacket := connections[0]
	firstPacket.TimeShift = 0
	for i, connection := range connections {
		connections[i].TimeShift = connection.Time.Sub(firstPacket.Time)
	}
	return connections, nil
}
