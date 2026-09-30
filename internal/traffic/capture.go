package traffic

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcapgo"
	pp "github.com/pires/go-proxyproto"
)

type Connection struct {
	TimeShift   time.Duration
	Time        time.Time
	SrcIP       string
	DstIP       string
	SrcPort     string
	DstPort     string
	Type        string
	Subtype     string
	Size        int
	Process     *ProcessInfo
	Capture     string
	SYN         bool
	ACK         bool
	RST         bool
	ClientIP    string
	Fingerprint string
}

type ProcessInfo struct {
	PID              *int64 `json:"pid,omitempty"`
	ProcessName      string `json:"process_name,omitempty"`
	BytesSent        int64  `json:"bytes_sent,omitempty"`
	BytesReceived    int64  `json:"bytes_received,omitempty"`
	GeoIPCountryCode string `json:"geoip_country_code,omitempty"`
	GeoIPCountryName string `json:"geoip_country_name,omitempty"`
	GeoIPASN         *int64 `json:"geoip_asn,omitempty"`
	GeoIPASOrg       string `json:"geoip_as_org,omitempty"`
	GeoIPCity        string `json:"geoip_city,omitempty"`
	GeoIPPostalCode  string `json:"geoip_postal_code,omitempty"`
	FirstSeen        *int64 `json:"first_seen,omitempty"`
	LastSeen         *int64 `json:"last_seen,omitempty"`
	matchFirstSeen   *time.Time
	matchLastSeen    *time.Time
}

type trafficConnectionSidecar struct {
	Protocol         string            `json:"protocol"`
	LocalAddr        string            `json:"local_addr"`
	RemoteAddr       string            `json:"remote_addr"`
	PID              any               `json:"pid"`
	ProcessName      string            `json:"process_name"`
	FirstSeen        trafficSystemTime `json:"first_seen"`
	LastSeen         trafficSystemTime `json:"last_seen"`
	BytesSent        int64             `json:"bytes_sent"`
	BytesReceived    int64             `json:"bytes_received"`
	GeoIPCountryCode string            `json:"geoip_country_code"`
	GeoIPCountryName string            `json:"geoip_country_name"`
	GeoIPASN         any               `json:"geoip_asn"`
	GeoIPASOrg       string            `json:"geoip_as_org"`
	GeoIPCity        string            `json:"geoip_city"`
	GeoIPPostalCode  string            `json:"geoip_postal_code"`
}

type trafficSystemTime struct {
	Secs  int64 `json:"secs_since_epoch"`
	Nanos int64 `json:"nanos_since_epoch"`
}

type trafficProcessLookup map[trafficProcessKey][]ProcessInfo

type trafficProcessKey struct {
	Protocol string
	Local    string
	Remote   string
}

func ReadPcapFile(ctx context.Context, path string) ([]Connection, []SourceIssue, error) {
	file, err := os.Stat(path)
	if err != nil {
		return nil, nil, err
	}
	if file.IsDir() {
		return nil, nil, fmt.Errorf("%s is a directory", path)
	}
	processLookup, issues, err := loadTrafficProcessLookup(ctx, path+".connections.jsonl")
	if err != nil {
		return nil, issues, err
	}
	var connections []Connection
	var firstPacketTime time.Time
	err = walkTrafficPackets(ctx, path, func(packet gopacket.Packet, _ layers.LinkType) error {
		connection, ok := extractTrafficConnection(packet, processLookup)
		if !ok {
			return nil
		}
		if firstPacketTime.IsZero() {
			firstPacketTime = packet.Metadata().Timestamp
		}
		connection.TimeShift = packet.Metadata().Timestamp.Sub(firstPacketTime)
		connection.Capture = filepath.Base(path)
		connections = append(connections, connection)
		return nil
	})
	return connections, issues, err
}

func EnrichPcap(ctx context.Context, pcapPath, jsonlPath, outputPath string) ([]SourceIssue, error) {
	processLookup, issues, err := loadTrafficProcessLookup(ctx, jsonlPath)
	if err != nil {
		return issues, err
	}
	return issues, enrichPcap(ctx, pcapPath, outputPath, processLookup)
}

func enrichPcap(ctx context.Context, pcapPath, outputPath string, processLookup trafficProcessLookup) error {
	file, err := os.CreateTemp(filepath.Dir(outputPath), ".traffic-enrich-*")
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}(file)

	var writer *pcapgo.NgWriter
	interfaces := make(map[layers.LinkType]int)
	err = walkTrafficPackets(ctx, pcapPath, func(packet gopacket.Packet, link layers.LinkType) error {
		if writer == nil {
			writer, err = pcapgo.NewNgWriter(file, link)
			if err != nil {
				return err
			}
			interfaces[link] = 0
		}
		index, ok := interfaces[link]
		if !ok {
			index, err = writer.AddInterface(pcapgo.NgInterface{LinkType: link, SnapLength: 0, TimestampResolution: 9})
			if err != nil {
				return err
			}
			interfaces[link] = index
		}
		connection, _ := extractTrafficConnection(packet, processLookup)
		options := pcapgo.NgPacketOptions{}
		if comment := buildTrafficProcessComment(connection.Process); comment != "" {
			options.Comments = []string{comment}
		}
		ci := packet.Metadata().CaptureInfo
		ci.InterfaceIndex = index
		return writer.WritePacketWithOptions(ci, packet.Data(), options)
	})
	if err != nil {
		return err
	}
	if writer == nil {
		writer, err = pcapgo.NewNgWriter(file, layers.LinkTypeEthernet)
	}
	if err != nil {
		return err
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), outputPath)
}

func extractTrafficConnection(packet gopacket.Packet, processLookup trafficProcessLookup) (Connection, bool) {
	connection := Connection{Size: packet.Metadata().CaptureLength, Time: packet.Metadata().Timestamp}
	if network := packet.NetworkLayer(); network != nil {
		connection.Fingerprint = packetFingerprint(append(append([]byte{}, network.LayerContents()...), network.LayerPayload()...))
	}
	src, dst, srcPort, dstPort, baseLayerIndex, ok := extractTrafficEndpoints(packet)
	if !ok {
		return Connection{}, false
	}
	connection.SrcIP = NormalizeTrafficIP(src)
	connection.DstIP = NormalizeTrafficIP(dst)
	if connection.SrcIP == "" || connection.DstIP == "" {
		return Connection{}, false
	}
	connection.SrcPort = srcPort
	connection.DstPort = dstPort
	connection.Type, connection.Subtype = extractTrafficProtocols(packet, baseLayerIndex)
	connection.Process = findTrafficProcess(processLookup, connection)
	if tcp, ok := packet.Layer(layers.LayerTypeTCP).(*layers.TCP); ok {
		connection.SYN, connection.ACK = tcp.SYN, tcp.ACK
		connection.RST = tcp.RST
	}
	if transport := packet.TransportLayer(); transport != nil {
		if header, readErr := pp.Read(bufio.NewReader(bytes.NewReader(transport.LayerPayload()))); readErr == nil && header.SourceAddr != nil && header.DestinationAddr != nil {
			srcIP, _, srcErr := net.SplitHostPort(header.SourceAddr.String())
			if srcErr == nil {
				// Preserve the captured network endpoints for internal topology.
				connection.ClientIP = NormalizeTrafficIP(srcIP)
			}
		}
	}
	if isIgnoredTrafficIP(connection.SrcIP) || isIgnoredTrafficIP(connection.DstIP) {
		return Connection{}, false
	}
	return connection, true
}

func isIgnoredTrafficIP(value string) bool {
	ip := net.ParseIP(strings.TrimSpace(value))
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	if ipv4 := ip.To4(); ipv4 != nil {
		if ipv4.Equal(net.IPv4bcast) {
			return true
		}
	}
	// RFC1918/ULA addresses remain available for target-network replay.
	return !ip.IsPrivate() && !IsPublicTrafficIP(value)
}

func buildTrafficProcessComment(process *ProcessInfo) string {
	if process == nil {
		return ""
	}
	parts := make([]string, 0, 8)
	if process.GeoIPCountryCode != "" {
		parts = append(parts, "Loc:"+process.GeoIPCountryCode)
	}
	if process.GeoIPCity != "" {
		parts = append(parts, "City:"+process.GeoIPCity)
	}
	if process.GeoIPPostalCode != "" {
		parts = append(parts, "ZIP:"+process.GeoIPPostalCode)
	}
	if process.GeoIPASN != nil {
		parts = append(parts, "AS"+strconv.FormatInt(*process.GeoIPASN, 10))
	}
	if process.PID != nil {
		parts = append(parts, "PID:"+strconv.FormatInt(*process.PID, 10))
	}
	if process.ProcessName != "" {
		parts = append(parts, "Process:"+process.ProcessName)
	}
	if process.BytesSent > 0 {
		parts = append(parts, "Sent:"+strconv.FormatInt(process.BytesSent, 10))
	}
	if process.BytesReceived > 0 {
		parts = append(parts, "Recv:"+strconv.FormatInt(process.BytesReceived, 10))
	}
	return strings.Join(parts, " ")
}

func extractTrafficEndpoints(packet gopacket.Packet) (string, string, string, string, int, bool) {
	layersL := packet.Layers()
	if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
		if arp, ok := arpLayer.(*layers.ARP); ok {
			src := formatTrafficARPAddress(arp.SourceProtAddress, arp.SourceHwAddress)
			dst := formatTrafficARPAddress(arp.DstProtAddress, arp.DstHwAddress)
			if src != "" && dst != "" {
				return src, dst, "", "", findTrafficLayerIndex(layersL, arp.LayerType()), true
			}
		}
	}
	if network := packet.NetworkLayer(); network != nil {
		src := formatTrafficEndpoint(network.NetworkFlow().Src())
		dst := formatTrafficEndpoint(network.NetworkFlow().Dst())
		if src != "" && dst != "" {
			srcPort, dstPort := extractTrafficPorts(packet)
			return src, dst, srcPort, dstPort, findTrafficLayerIndex(layersL, network.LayerType()), true
		}
	}
	if link := packet.LinkLayer(); link != nil {
		src := formatTrafficEndpoint(link.LinkFlow().Src())
		dst := formatTrafficEndpoint(link.LinkFlow().Dst())
		if src != "" && dst != "" {
			return src, dst, "", "", findTrafficLayerIndex(layersL, link.LayerType()), true
		}
	}
	return "", "", "", "", -1, false
}

func extractTrafficPorts(packet gopacket.Packet) (string, string) {
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		if tcp, ok := tcpLayer.(*layers.TCP); ok {
			return strconv.Itoa(int(tcp.SrcPort)), strconv.Itoa(int(tcp.DstPort))
		}
	}
	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		if udp, ok := udpLayer.(*layers.UDP); ok {
			return strconv.Itoa(int(udp.SrcPort)), strconv.Itoa(int(udp.DstPort))
		}
	}
	return "", ""
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

func loadTrafficProcessLookup(ctx context.Context, path string) (trafficProcessLookup, []SourceIssue, error) {
	lookup := make(trafficProcessLookup)
	issues := make([]SourceIssue, 0)
	warn := func(err error) {
		issues = append(issues, SourceIssue{File: filepath.Base(path), Phase: "process_metadata", Error: err.Error()})
	}
	if err := ctx.Err(); err != nil {
		return lookup, issues, err
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return lookup, issues, nil
		}
		warn(err)
		return lookup, issues, nil
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	reader := bufio.NewReader(file)
	lineNumber := 0
	for {
		if err = ctx.Err(); err != nil {
			return lookup, issues, err
		}
		line, readErr := reader.ReadString('\n')
		lineNumber++
		line = strings.TrimSpace(line)
		if line != "" {
			var sidecar trafficConnectionSidecar
			if err = json.Unmarshal([]byte(line), &sidecar); err != nil {
				warn(fmt.Errorf("line %d: %w", lineNumber, err))
			} else {
				info := sidecar.toProcessInfo()
				protocol := strings.ToUpper(strings.TrimSpace(sidecar.Protocol))
				local := strings.TrimSpace(sidecar.LocalAddr)
				remote := strings.TrimSpace(sidecar.RemoteAddr)
				if protocol != "" && local != "" && remote != "" {
					for _, key := range []trafficProcessKey{
						{Protocol: protocol, Local: local, Remote: remote},
						{Protocol: protocol, Local: remote, Remote: local},
					} {
						lookup[key] = append(lookup[key], info)
					}
				}
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			warn(readErr)
			break
		}
	}
	return lookup, issues, nil
}

func (sidecar trafficConnectionSidecar) toProcessInfo() ProcessInfo {
	firstSeen := sidecar.FirstSeen.toTime()
	lastSeen := sidecar.LastSeen.toTime()
	return ProcessInfo{
		PID:              parseTrafficInt(sidecar.PID),
		ProcessName:      sidecar.ProcessName,
		BytesSent:        sidecar.BytesSent,
		BytesReceived:    sidecar.BytesReceived,
		GeoIPCountryCode: sidecar.GeoIPCountryCode,
		GeoIPCountryName: sidecar.GeoIPCountryName,
		GeoIPASN:         parseTrafficInt(sidecar.GeoIPASN),
		GeoIPASOrg:       sidecar.GeoIPASOrg,
		GeoIPCity:        sidecar.GeoIPCity,
		GeoIPPostalCode:  sidecar.GeoIPPostalCode,
		FirstSeen:        timePtrToUnixNano(firstSeen),
		LastSeen:         timePtrToUnixNano(lastSeen),
		matchFirstSeen:   firstSeen,
		matchLastSeen:    lastSeen,
	}
}

func (st trafficSystemTime) toTime() *time.Time {
	if st.Secs == 0 && st.Nanos == 0 {
		return nil
	}
	return new(time.Unix(st.Secs, st.Nanos))
}

func timePtrToUnixNano(value *time.Time) *int64 {
	if value == nil {
		return nil
	}
	return new(value.UnixNano())
}

func parseTrafficInt(value any) *int64 {
	switch typed := value.(type) {
	case nil:
		return nil
	case float64:
		return new(int64(typed))
	case int64:
		return new(typed)
	case int:
		return new(int64(typed))
	case string:
		if typed == "" {
			return nil
		}
		parsed, err := strconv.ParseInt(typed, 10, 64)
		if err != nil {
			return nil
		}
		return &parsed
	default:
		return nil
	}
}

func findTrafficProcess(lookup trafficProcessLookup, connection Connection) *ProcessInfo {
	if len(lookup) == 0 {
		return nil
	}
	protocol := strings.ToUpper(strings.TrimSpace(connection.Type))
	if protocol == "IPV4" || protocol == "IPV6" || protocol == "IP" {
		protocol = strings.ToUpper(strings.TrimSpace(connection.Subtype))
	}
	if protocol == "" {
		return nil
	}
	keys := []trafficProcessKey{
		{
			Protocol: protocol,
			Local:    formatTrafficAddrPort(connection.SrcIP, connection.SrcPort),
			Remote:   formatTrafficAddrPort(connection.DstIP, connection.DstPort),
		},
	}
	if connection.SrcPort == "" || connection.DstPort == "" {
		keys = append(keys, trafficProcessKey{Protocol: protocol, Local: connection.SrcIP, Remote: connection.DstIP})
	}
	for _, key := range keys {
		if match := findTrafficProcessByKey(lookup, key, connection.Time); match != nil {
			return match
		}
	}
	return nil
}

func findTrafficProcessByKey(lookup trafficProcessLookup, key trafficProcessKey, packetTime time.Time) *ProcessInfo {
	matches := lookup[key]
	if len(matches) == 0 {
		return nil
	}
	const slack = 2 * time.Second
	var fallback *ProcessInfo
	var best *ProcessInfo
	bestScore := time.Duration(1<<63 - 1)
	for i := range matches {
		match := &matches[i]
		if match.matchFirstSeen == nil || match.matchLastSeen == nil {
			if fallback == nil {
				fallback = match
			}
			continue
		}
		first := match.matchFirstSeen.Add(-slack)
		last := match.matchLastSeen.Add(slack)
		if packetTime.Before(first) || packetTime.After(last) {
			continue
		}
		score := time.Duration(0)
		if packetTime.Before(*match.matchFirstSeen) {
			score = match.matchFirstSeen.Sub(packetTime)
		} else if packetTime.After(*match.matchLastSeen) {
			score = packetTime.Sub(*match.matchLastSeen)
		}
		if score < bestScore {
			bestScore = score
			best = match
		}
	}
	if best != nil {
		return best
	}
	return fallback
}

func formatTrafficAddrPort(ip, port string) string {
	if port == "" {
		return ip
	}
	return net.JoinHostPort(ip, port)
}

func EnrichPcapDirWithContext(ctx context.Context, path string) ([]string, []SourceIssue, error) {
	d, err := os.Stat(path)
	if err != nil {
		return nil, nil, err
	}
	if !d.IsDir() {
		return nil, nil, fmt.Errorf("%s is a file", path)
	}
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}
	generated := make([]string, 0)
	allIssues := make([]SourceIssue, 0)
	for _, file := range dir {
		if err = ctx.Err(); err != nil {
			return generated, allIssues, err
		}
		if file.IsDir() || !isOriginalTrafficCapture(file.Name()) {
			continue
		}
		pcapPath := filepath.Join(path, file.Name())
		jsonl := filepath.Join(path, file.Name()+".connections.jsonl")
		output := filepath.Join(path, file.Name()+".enrich.pcap")
		issues, enrichErr := EnrichPcap(ctx, pcapPath, jsonl, output)
		allIssues = append(allIssues, issues...)
		if err = enrichErr; err != nil {
			allIssues = append(allIssues, SourceIssue{File: file.Name(), Phase: "enrichment", Error: err.Error()})
		} else {
			generated = append(generated, output)
		}
	}
	return generated, allIssues, ctx.Err()
}

// frpcPcapName 是 frpc pod capture sidecar 写入的固定文件名。
const frpcPcapName = "frpc.pcap"

// PcapDirResult 读取靶机流量目录的结果：
//   - Connections：普通 pod 流量，用于拓扑展示。
//   - FrpcIPs：frpc pod 经 Proxy Protocol 传递的真实客户端 IP。
type PcapDirResult struct {
	Connections  []Connection
	FrpcIPs      []string
	Accesses     []Access
	SourceIssues []SourceIssue
	Files        []CaptureFileResult
}

func ReadPcapDir(ctx context.Context, path string, proxyPorts map[uint16]bool) (PcapDirResult, error) {
	if err := ctx.Err(); err != nil {
		return PcapDirResult{}, err
	}
	d, err := os.Stat(path)
	if err != nil {
		return PcapDirResult{}, err
	}
	if !d.IsDir() {
		return PcapDirResult{}, fmt.Errorf("%s is a file", path)
	}
	dir, err := os.ReadDir(path)
	if err != nil {
		return PcapDirResult{}, err
	}

	connections := make([]Connection, 0)
	frpcIPSet := make(map[string]struct{})
	accesses := make([]Access, 0)
	issues := make([]SourceIssue, 0)
	files := make([]CaptureFileResult, 0)

	for _, file := range dir {
		if err = ctx.Err(); err != nil {
			return PcapDirResult{}, err
		}
		if file.IsDir() || !isOriginalTrafficCapture(file.Name()) {
			continue
		}
		fullPath := filepath.Join(path, file.Name())
		if file.Name() == frpcPcapName {
			// frpc 流量：只提取 Proxy Protocol 中的真实客户端 IP，其余包跳过。
			observations, readErr := extractFrpcProxyAccesses(ctx, fullPath, proxyPorts)
			if ctx.Err() != nil {
				return PcapDirResult{}, ctx.Err()
			}
			if readErr != nil {
				issues = append(issues, SourceIssue{File: file.Name(), Phase: "proxy", Error: readErr.Error()})
			}
			accesses = append(accesses, observations...)
			for _, access := range observations {
				frpcIPSet[access.IP] = struct{}{}
			}
		} else {
			// 普通 pod 流量：完整分析，进入拓扑展示。
			packetConnections, metadataIssues, readErr := ReadPcapFile(ctx, fullPath)
			issues = append(issues, metadataIssues...)
			if ctx.Err() != nil {
				return PcapDirResult{}, ctx.Err()
			}
			status := "processed"
			if readErr != nil {
				status = "partial"
				issues = append(issues, SourceIssue{File: file.Name(), Phase: "replay", Error: readErr.Error()})
			}
			files = append(files, CaptureFileResult{File: file.Name(), Status: status, Packets: int64(len(packetConnections))})
			connections = append(connections, packetConnections...)
		}
	}

	if len(connections) > 0 {
		slices.SortStableFunc(connections, func(c1 Connection, c2 Connection) int { return c1.Time.Compare(c2.Time) })
		connections = deduplicateTraffic(connections)
		firstPacket := connections[0]
		for i, connection := range connections {
			connections[i].TimeShift = connection.Time.Sub(firstPacket.Time)
		}
	}

	frpcIPs := make([]string, 0, len(frpcIPSet))
	for ip := range frpcIPSet {
		frpcIPs = append(frpcIPs, ip)
	}
	slices.Sort(frpcIPs)

	return PcapDirResult{
		Connections:  connections,
		FrpcIPs:      frpcIPs,
		Accesses:     accesses,
		SourceIssues: issues,
		Files:        files,
	}, nil
}

// extractFrpcProxyAccesses reads trusted FRPC-to-nginx client evidence.
func extractFrpcProxyAccesses(ctx context.Context, path string, proxyPorts map[uint16]bool) ([]Access, error) {
	accesses := make([]Access, 0)
	warnings := make(map[string]bool)
	consume := func(data []byte, evidence Evidence) {
		header, err := pp.Read(bufio.NewReader(bytes.NewReader(data)))
		if err != nil {
			if bytes.HasPrefix(data, []byte("PROXY ")) || bytes.HasPrefix(data, []byte("\r\n\r\n")) {
				warnings["proxy_header_incomplete"] = true
			}
			return
		}
		if header.SourceAddr == nil {
			return
		}
		ip, _, err := net.SplitHostPort(header.SourceAddr.String())
		if err == nil {
			accesses = append(accesses, Access{IP: NormalizeTrafficIP(ip), Time: evidence.Time, Source: "proxy_protocol", Capture: evidence.Capture})
		}
	}
	streams := newStreamCollector(consume, func(code string) { warnings[code] = true })
	streams.prefixOnly = true
	err := walkTrafficPackets(ctx, path, func(packet gopacket.Packet, _ layers.LinkType) error {
		// FRPC delivers PROXY v2 to nginx over loopback ports 10000+.
		// Filter the claimed client IP later, not this local transport envelope.
		src, dst, srcPort, dstPort, _, ok := extractTrafficEndpoints(packet)
		port, _ := strconv.Atoi(dstPort)
		ip := net.ParseIP(dst)
		if !ok || ip == nil || !ip.IsLoopback() || port < 0 || port > 65535 || !proxyPorts[uint16(port)] {
			return nil
		}
		c := Connection{SrcIP: NormalizeTrafficIP(src), DstIP: NormalizeTrafficIP(dst), SrcPort: srcPort, DstPort: dstPort, Type: "TCP", Time: packet.Metadata().Timestamp, Capture: filepath.Base(path)}
		e := packetEvidence(c)
		if tcp, ok := packet.Layer(layers.LayerTypeTCP).(*layers.TCP); ok {
			streams.add(tcp, e)
		} else if transport := packet.TransportLayer(); transport != nil {
			consume(transport.LayerPayload(), e)
		}
		return nil
	})
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	streams.flush()
	// Never silently publish a complete-looking anti-cheat evidence index after
	// an analysis resource limit or a capture gap.
	if len(warnings) > 0 {
		codes := make([]string, 0, len(warnings))
		for code := range warnings {
			codes = append(codes, code)
		}
		slices.Sort(codes)
		return accesses, fmt.Errorf("incomplete FRP TCP evidence in %s: %s (read error: %v)", path, strings.Join(codes, ","), err)
	}
	return accesses, err
}
