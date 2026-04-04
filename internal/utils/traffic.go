package utils

import (
	"CBCTF/internal/log"
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
		src, dst, baseLayerIndex, ok := extractTrafficEndpoints(packet)
		if !ok {
			continue
		}
		connection.SrcIP = src
		connection.DstIP = dst
		connection.Type, connection.Subtype = extractTrafficProtocols(packet, baseLayerIndex)
		if transport := packet.TransportLayer(); transport != nil {
			if header, readErr := pp.Read(bufio.NewReader(bytes.NewReader(transport.LayerPayload()))); readErr == nil {
				srcIP, _, srcErr := net.SplitHostPort(header.SourceAddr.String())
				dstIP, _, dstErr := net.SplitHostPort(header.DestinationAddr.String())
				if srcErr == nil && dstErr == nil {
					connection.SrcIP = srcIP
					connection.DstIP = dstIP
					connection.Subtype = "Proxy"
				}
			}
		}
		connections = append(connections, connection)
	}
	return connections, nil
}

func extractTrafficEndpoints(packet gopacket.Packet) (string, string, int, bool) {
	layersL := packet.Layers()
	if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
		if arp, ok := arpLayer.(*layers.ARP); ok {
			src := formatTrafficARPAddress(arp.SourceProtAddress, arp.SourceHwAddress)
			dst := formatTrafficARPAddress(arp.DstProtAddress, arp.DstHwAddress)
			if src != "" && dst != "" {
				return src, dst, findTrafficLayerIndex(layersL, arp.LayerType()), true
			}
		}
	}
	if network := packet.NetworkLayer(); network != nil {
		src := formatTrafficEndpoint(network.NetworkFlow().Src())
		dst := formatTrafficEndpoint(network.NetworkFlow().Dst())
		if src != "" && dst != "" {
			return src, dst, findTrafficLayerIndex(layersL, network.LayerType()), true
		}
	}
	if link := packet.LinkLayer(); link != nil {
		src := formatTrafficEndpoint(link.LinkFlow().Src())
		dst := formatTrafficEndpoint(link.LinkFlow().Dst())
		if src != "" && dst != "" {
			return src, dst, findTrafficLayerIndex(layersL, link.LayerType()), true
		}
	}
	return "", "", -1, false
}

func formatTrafficARPAddress(protocolAddress, hardwareAddress []byte) string {
	switch len(protocolAddress) {
	case net.IPv4len, net.IPv6len:
		return net.IP(protocolAddress).String()
	}
	if len(protocolAddress) > 0 {
		return fmt.Sprintf("%x", protocolAddress)
	}
	if len(hardwareAddress) > 0 {
		return net.HardwareAddr(hardwareAddress).String()
	}
	return ""
}

func formatTrafficEndpoint(endpoint gopacket.Endpoint) string {
	value := strings.TrimSpace(endpoint.String())
	if value == "" {
		return ""
	}
	return value
}

func extractTrafficProtocols(packet gopacket.Packet, baseLayerIndex int) (string, string) {
	layersL := packet.Layers()
	if baseLayerIndex < 0 || baseLayerIndex >= len(layersL) {
		baseLayerIndex = -1
	}

	protocols := make([]string, 0, 2)
	for i := baseLayerIndex + 1; i < len(layersL); i++ {
		layerType := normalizeTrafficLayerType(layersL[i].LayerType())
		if layerType == "" {
			continue
		}
		protocols = append(protocols, layerType)
		if len(protocols) >= 2 {
			break
		}
	}

	if len(protocols) == 0 && baseLayerIndex >= 0 {
		protocols = append(protocols, normalizeTrafficLayerType(layersL[baseLayerIndex].LayerType()))
	}

	protocol := ""
	subtype := ""
	if len(protocols) > 0 {
		protocol = protocols[0]
	}
	if len(protocols) > 1 {
		subtype = protocols[1]
	}
	if subtype == "" {
		if application := packet.ApplicationLayer(); application != nil {
			subtype = normalizeTrafficLayerType(application.LayerType())
		}
	}
	return protocol, subtype
}

func normalizeTrafficLayerType(layerType gopacket.LayerType) string {
	switch layerType {
	case gopacket.LayerTypePayload, gopacket.LayerTypeDecodeFailure:
		return ""
	}
	name := strings.TrimSpace(layerType.String())
	if name == "" {
		return ""
	}
	return strings.TrimPrefix(name, "LayerType")
}

func findTrafficLayerIndex(layersL []gopacket.Layer, target gopacket.LayerType) int {
	for i, layer := range layersL {
		if layer.LayerType() == target {
			return i
		}
	}
	return -1
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
		packetConnections, readErr := ReadPcapFile(fmt.Sprintf("%s/%s", path, file.Name()))
		if readErr != nil {
			log.Logger.Warningf("Failed to read pcap file %s: %s", file.Name(), readErr.Error())
			continue
		}
		connections = append(connections, packetConnections...)
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
