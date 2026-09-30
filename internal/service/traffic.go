package service

import (
	"context"
	"fmt"
	"net/netip"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"
	"CBCTF/internal/traffic"
)

type trafficNodeAggregate struct {
	ID            string
	Label         string
	IP            string
	Kind          string
	Side          string
	Zone          string
	Service       string
	Services      []string
	Bytes         int64
	Packets       int64
	Connections   int64
	protocolBytes map[string]int64
	processes     map[string]*trafficProcessAggregate
	flows         map[string]bool
}

type trafficEdgeAggregate struct {
	ID            string
	Source        string
	Target        string
	Direction     string
	Kind          string
	SourceService string
	TargetService string
	Bytes         int64
	Packets       int64
	Connections   int64
	protocolBytes map[string]int64
	appBytes      map[string]int64
	processes     map[string]*trafficProcessAggregate
	flows         map[string]bool
}

type trafficProcessAggregate struct {
	Info    traffic.ProcessInfo
	Bytes   int64
	Packets int64
}

type trafficBucketAggregate struct {
	TimestampMs  int64
	Bytes        int64
	Packets      int64
	IngressBytes int64
	EgressBytes  int64
}

type trafficRankingAggregate struct {
	Label         string
	IP            string
	Bytes         int64
	Packets       int64
	Connections   int64
	DominantProto string
	DominantApp   string
	DominantProc  string
	Direction     string
	Processes     []resp.TrafficProcessResp
}

const DefaultTrafficDurationMs int64 = 1000

func GetTraffic(ctx context.Context, victim model.Victim, form dto.GetTrafficForm) (resp.TrafficTopologyResp, model.RetVal) {
	snapshot, ret := loadTrafficSnapshot(ctx, victim)
	if !ret.OK {
		return resp.TrafficTopologyResp{}, ret
	}
	connections := snapshot.Connections

	// 从 DB 读取该靶机所有已知 IP（不受时间窗口限制）
	ips, ret := db.InitTrafficRepo(db.DB.WithContext(ctx)).GetVictimIPs(victim.ID)
	if !ret.OK {
		return resp.TrafficTopologyResp{}, ret
	}

	if len(connections) == 0 {
		empty := emptyTrafficTopology(victim, form)
		empty.IPs = ips
		empty.SourceIssues = snapshot.SourceIssues
		empty.Files = snapshot.Files
		return empty, model.SuccessRetVal()
	}

	totalDuration := calcTrafficTotalDuration(connections)
	start, end := clampTrafficWindow(form.TimeShift, form.Duration, totalDuration)
	internalIPs := victim.TrafficInternalIPs()
	servicesByIP := collectVictimServicesByIP(victim)
	bucketSize := max(DefaultTrafficDurationMs, (totalDuration+1999)/2000)
	timelineBuckets := buildTrafficTimelineBuckets(connections, internalIPs, bucketSize)
	windowConnections := sliceTrafficConnections(connections, start, end)
	totalPackets := int64(len(windowConnections))

	nodes := make(map[string]*trafficNodeAggregate)
	edges := make(map[string]*trafficEdgeAggregate)
	// Keep configured target nodes present during idle replay windows.
	for ip := range internalIPs {
		services := servicesByIP[ip]
		nodes[ip] = &trafficNodeAggregate{ID: ip, IP: ip, Label: buildTrafficNodeLabel(ip, true, services),
			Kind: "victim", Side: "center", Zone: "victim", Service: dominantTrafficService(services), Services: services,
			protocolBytes: make(map[string]int64), processes: make(map[string]*trafficProcessAggregate), flows: make(map[string]bool)}
	}

	summary := resp.TrafficSummaryResp{}
	allFlows := make(map[string]bool)
	topTalkers := make([]trafficRankingAggregate, 0)
	topEdges := make([]trafficRankingAggregate, 0)

	for _, connection := range windowConnections {
		packetBytes := int64(connection.Size)
		summary.TotalBytes += packetBytes
		summary.TotalPackets++

		srcInternal := internalIPs[connection.SrcIP]
		dstInternal := internalIPs[connection.DstIP]

		direction := trafficDirection(srcInternal, dstInternal)
		if direction == "ingress" {
			summary.IngressBytes += packetBytes
		} else if direction == "egress" {
			summary.EgressBytes += packetBytes
		} else if direction == "internal" {
			summary.InternalBytes += packetBytes
		} else {
			summary.ExternalBytes += packetBytes
		}

		protocol := normalizeTrafficProtocol(connection.Type)
		app := normalizeTrafficSubtype(connection.Subtype)

		edgeID := buildTrafficEdgeID(connection.SrcIP, connection.DstIP, direction)
		edge := edges[edgeID]
		if edge == nil {
			edge = &trafficEdgeAggregate{
				ID:            edgeID,
				Source:        connection.SrcIP,
				Target:        connection.DstIP,
				Direction:     direction,
				Kind:          trafficEdgeKind(srcInternal, dstInternal),
				SourceService: dominantTrafficService(servicesByIP[connection.SrcIP]),
				TargetService: dominantTrafficService(servicesByIP[connection.DstIP]),
				protocolBytes: make(map[string]int64),
				appBytes:      make(map[string]int64),
				processes:     make(map[string]*trafficProcessAggregate),
				flows:         make(map[string]bool),
			}
			edges[edgeID] = edge
		}
		edge.Bytes += packetBytes
		edge.Packets++
		flow := connection.FlowKey()
		allFlows[flow] = true
		edge.flows[flow] = true
		edge.Connections = int64(len(edge.flows))
		edge.protocolBytes[protocol] += packetBytes
		edge.appBytes[app] += packetBytes
		addTrafficProcess(edge.processes, connection.Process, packetBytes)

		for _, ip := range []string{connection.SrcIP, connection.DstIP} {
			services := servicesByIP[ip]
			node := nodes[ip]
			if node == nil {
				node = &trafficNodeAggregate{
					ID:            ip,
					Label:         buildTrafficNodeLabel(ip, internalIPs[ip], services),
					IP:            ip,
					Kind:          trafficNodeKind(internalIPs[ip]),
					Side:          trafficNodeSide(ip, srcInternal, dstInternal, internalIPs),
					Zone:          trafficNodeZone(ip, internalIPs),
					Service:       dominantTrafficService(services),
					Services:      services,
					protocolBytes: make(map[string]int64),
					processes:     make(map[string]*trafficProcessAggregate),
					flows:         make(map[string]bool),
				}
				nodes[ip] = node
			}
			node.Bytes += packetBytes
			node.Packets++
			node.flows[flow] = true
			node.Connections = int64(len(node.flows))
			node.protocolBytes[protocol] += packetBytes
			addTrafficProcess(node.processes, connection.Process, packetBytes)
		}
	}

	summary.PeakTimeMs, summary.PeakBytes = computeTrafficPeak(timelineBuckets)

	nodeList := make([]resp.TrafficNodeResp, 0, len(nodes))
	internalNodeCount := 0
	externalNodeCount := 0
	for _, node := range nodes {
		protocols := sortTrafficProtocolKeys(node.protocolBytes)
		processes := buildTrafficProcessResp(node.processes, 4)
		dominantProc := dominantTrafficProcess(processes)
		nodeList = append(nodeList, resp.TrafficNodeResp{
			ID:            node.ID,
			Label:         node.Label,
			IP:            node.IP,
			Kind:          node.Kind,
			Side:          node.Side,
			Zone:          node.Zone,
			Service:       node.Service,
			Services:      node.Services,
			Bytes:         node.Bytes,
			Packets:       node.Packets,
			Connections:   node.Connections,
			Protocols:     protocols,
			DominantProto: dominantTrafficKey(node.protocolBytes),
			DominantProc:  dominantProc,
			Processes:     processes,
		})
		if node.Kind == "victim" {
			internalNodeCount++
		} else {
			externalNodeCount++
			topTalkers = append(topTalkers, trafficRankingAggregate{
				Label:         node.Label,
				IP:            node.IP,
				Bytes:         node.Bytes,
				Packets:       node.Packets,
				Connections:   node.Connections,
				DominantProto: dominantTrafficKey(node.protocolBytes),
				DominantProc:  dominantProc,
				Processes:     processes,
			})
		}
	}
	sort.Slice(nodeList, func(i, j int) bool {
		if nodeList[i].Kind != nodeList[j].Kind {
			return nodeList[i].Kind < nodeList[j].Kind
		}
		if nodeList[i].Bytes != nodeList[j].Bytes {
			return nodeList[i].Bytes > nodeList[j].Bytes
		}
		return nodeList[i].IP < nodeList[j].IP
	})

	maxEdgeBytes := int64(1)
	for _, edge := range edges {
		if edge.Bytes > maxEdgeBytes {
			maxEdgeBytes = edge.Bytes
		}
	}

	edgeList := make([]resp.TrafficEdgeResp, 0, len(edges))
	for _, edge := range edges {
		protocols := sortTrafficProtocolKeys(edge.protocolBytes)
		dominantProto := dominantTrafficKey(edge.protocolBytes)
		dominantApp := dominantTrafficKey(edge.appBytes)
		processes := buildTrafficProcessResp(edge.processes, 4)
		dominantProc := dominantTrafficProcess(processes)
		edgeList = append(edgeList, resp.TrafficEdgeResp{
			ID:            edge.ID,
			Source:        edge.Source,
			Target:        edge.Target,
			Direction:     edge.Direction,
			Kind:          edge.Kind,
			SourceService: edge.SourceService,
			TargetService: edge.TargetService,
			Bytes:         edge.Bytes,
			Packets:       edge.Packets,
			Connections:   edge.Connections,
			Weight:        float64(edge.Bytes) / float64(maxEdgeBytes),
			Intensity:     trafficIntensity(edge.Bytes, maxEdgeBytes),
			Protocols:     protocols,
			DominantProto: dominantProto,
			DominantApp:   dominantApp,
			DominantProc:  dominantProc,
			Processes:     processes,
		})
		topEdges = append(topEdges, trafficRankingAggregate{
			Label:         fmt.Sprintf("%s -> %s", edge.Source, edge.Target),
			IP:            edge.ID,
			Bytes:         edge.Bytes,
			Packets:       edge.Packets,
			Connections:   edge.Connections,
			DominantProto: dominantProto,
			DominantApp:   dominantApp,
			DominantProc:  dominantProc,
			Direction:     edge.Direction,
			Processes:     processes,
		})
	}
	sort.Slice(edgeList, func(i, j int) bool {
		if edgeList[i].Bytes != edgeList[j].Bytes {
			return edgeList[i].Bytes > edgeList[j].Bytes
		}
		return edgeList[i].ID < edgeList[j].ID
	})

	timeline := make([]resp.TrafficTimelineBucketResp, 0, len(timelineBuckets))
	timestamps := make([]int64, 0, len(timelineBuckets))
	for timestampMs := range timelineBuckets {
		timestamps = append(timestamps, timestampMs)
	}
	slices.Sort(timestamps)
	for _, timestampMs := range timestamps {
		bucket := timelineBuckets[timestampMs]
		timeline = append(timeline, resp.TrafficTimelineBucketResp{
			TimestampMs:  bucket.TimestampMs,
			Bytes:        bucket.Bytes,
			Packets:      bucket.Packets,
			IngressBytes: bucket.IngressBytes,
			EgressBytes:  bucket.EgressBytes,
		})
	}

	sort.Slice(topTalkers, func(i, j int) bool {
		if topTalkers[i].Bytes != topTalkers[j].Bytes {
			return topTalkers[i].Bytes > topTalkers[j].Bytes
		}
		return topTalkers[i].IP < topTalkers[j].IP
	})
	sort.Slice(topEdges, func(i, j int) bool {
		if topEdges[i].Bytes != topEdges[j].Bytes {
			return topEdges[i].Bytes > topEdges[j].Bytes
		}
		return topEdges[i].Label < topEdges[j].Label
	})

	summary.InternalNodes = internalNodeCount
	summary.ExternalNodes = externalNodeCount
	summary.VisibleEdges = len(edgeList)
	summary.VisibleNodes = len(nodeList)
	summary.ProcessCount = countTrafficProcesses(windowConnections)
	summary.TotalConnections = len(allFlows)

	return resp.TrafficTopologyResp{
		StartedAt:        connections[0].Time.UTC().Format(time.RFC3339Nano),
		SourceIssues:     snapshot.SourceIssues,
		Files:            snapshot.Files,
		TimelineBucketMs: bucketSize,
		Window: resp.TrafficWindowResp{
			Start:      start,
			End:        end,
			Duration:   end - start,
			Total:      totalDuration,
			TotalCount: totalPackets,
		},
		TotalDuration:   totalDuration,
		AvailableSlices: availableTrafficSlices(totalDuration),
		Center: resp.TrafficCenterResp{
			Label:   buildTrafficCenterLabel(victim),
			Exposed: victim.RemoteAddr(),
		},
		Summary:    summary,
		Nodes:      nodeList,
		Edges:      edgeList,
		Timeline:   timeline,
		TopTalkers: buildTrafficRankingResp(topTalkers, 6),
		TopEdges:   buildTrafficRankingResp(topEdges, 6),
		IPs:        ips,
	}, model.SuccessRetVal()
}

func loadTrafficSnapshot(ctx context.Context, victim model.Victim) (*traffic.PcapDirResult, model.RetVal) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	snapshot, ret := redis.GetTraffic(ctx, victim)
	if ret.OK && snapshot != nil {
		return snapshot, ret
	}
	result, err := traffic.ReadPcapDir(ctx, victim.TrafficBasePath(), victim.TrafficProxyPorts())
	if err != nil {
		return nil, model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}
	if ret = redis.StoreTraffic(ctx, victim, &result); !ret.OK {
		log.Logger.Warningf("Traffic cache unavailable: victim_id=%d reason=%s", victim.ID, ret.Msg)
	}
	return &result, model.SuccessRetVal()
}

func emptyTrafficTopology(victim model.Victim, form dto.GetTrafficForm) resp.TrafficTopologyResp {
	duration := form.Duration
	if duration <= 0 {
		duration = DefaultTrafficDurationMs
	}
	nodes := make([]resp.TrafficNodeResp, 0)
	services := collectVictimServicesByIP(victim)
	for ip := range victim.TrafficInternalIPs() {
		nodes = append(nodes, resp.TrafficNodeResp{ID: ip, IP: ip, Label: buildTrafficNodeLabel(ip, true, services[ip]),
			Kind: "victim", Side: "center", Zone: "victim", Service: dominantTrafficService(services[ip]), Services: services[ip], Protocols: []string{}})
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].IP < nodes[j].IP })
	return resp.TrafficTopologyResp{
		Window: resp.TrafficWindowResp{
			Start:    form.TimeShift,
			End:      form.TimeShift + duration,
			Duration: duration,
			Total:    0,
		},
		TotalDuration:   0,
		AvailableSlices: []int64{1000, 5000, 15000, 30000, 60000},
		Center: resp.TrafficCenterResp{
			Label:   buildTrafficCenterLabel(victim),
			Exposed: victim.RemoteAddr(),
		},
		Summary:    resp.TrafficSummaryResp{InternalNodes: len(nodes), VisibleNodes: len(nodes)},
		Nodes:      nodes,
		Edges:      make([]resp.TrafficEdgeResp, 0),
		Timeline:   make([]resp.TrafficTimelineBucketResp, 0),
		TopTalkers: make([]resp.TrafficRankingResp, 0),
		TopEdges:   make([]resp.TrafficRankingResp, 0),
	}
}

func calcTrafficTotalDuration(connections []traffic.Connection) int64 {
	if len(connections) == 0 {
		return 0
	}
	durationMs := connections[len(connections)-1].Time.Sub(connections[0].Time).Milliseconds()
	return durationMs + 1
}

func clampTrafficWindow(start, duration, total int64) (int64, int64) {
	if start < 0 {
		start = 0
	}
	if duration <= 0 {
		duration = DefaultTrafficDurationMs
	}
	if total > 0 && start > total {
		start = total
	}
	end := start + min(duration, max(0, total-start))
	return start, end
}

func sliceTrafficConnections(connections []traffic.Connection, start, end int64) []traffic.Connection {
	if len(connections) == 0 {
		return make([]traffic.Connection, 0)
	}
	first := sort.Search(len(connections), func(i int) bool { return connections[i].TimeShift.Milliseconds() >= start })
	last := sort.Search(len(connections), func(i int) bool { return connections[i].TimeShift.Milliseconds() >= end })
	return connections[first:last]
}

func buildTrafficTimelineBuckets(
	connections []traffic.Connection,
	internalIPs map[string]bool,
	bucketSizeMs int64,
) map[int64]*trafficBucketAggregate {
	if bucketSizeMs <= 0 {
		bucketSizeMs = DefaultTrafficDurationMs
	}
	buckets := make(map[int64]*trafficBucketAggregate)
	for _, connection := range connections {
		packetBytes := int64(connection.Size)
		timestampMs := connection.TimeShift.Milliseconds()
		bucketStart := (timestampMs / bucketSizeMs) * bucketSizeMs
		bucket := buckets[bucketStart]
		if bucket == nil {
			bucket = &trafficBucketAggregate{TimestampMs: bucketStart}
			buckets[bucketStart] = bucket
		}
		bucket.Bytes += packetBytes
		bucket.Packets++

		srcInternal := internalIPs[connection.SrcIP]
		dstInternal := internalIPs[connection.DstIP]
		switch trafficDirection(srcInternal, dstInternal) {
		case "ingress":
			bucket.IngressBytes += packetBytes
		case "egress":
			bucket.EgressBytes += packetBytes
		}
	}
	return buckets
}

func collectVictimServicesByIP(victim model.Victim) map[string][]string {
	servicesByIP := make(map[string][]string)
	for _, pod := range victim.Pods {
		addPodServicesByIP(servicesByIP, pod.Spec)
	}
	for _, podSpec := range victim.Spec.Pods {
		addPodServicesByIP(servicesByIP, podSpec)
	}
	for ip, services := range servicesByIP {
		servicesByIP[ip] = uniqueSortedTrafficServices(services)
	}
	return servicesByIP
}

func addPodServicesByIP(servicesByIP map[string][]string, podSpec model.PodSpec) {
	services := make([]string, 0, len(podSpec.Containers))
	for _, container := range podSpec.Containers {
		name := strings.TrimSpace(container.Name)
		if name == "" || isTrafficInfraContainer(name) {
			continue
		}
		services = append(services, name)
	}
	if len(services) == 0 {
		return
	}
	for _, network := range podSpec.Networks {
		ip := network.Attachment.IP
		if prefix, err := netip.ParsePrefix(ip); err == nil {
			ip = prefix.Addr().String()
		}
		if ip = traffic.NormalizeTrafficIP(ip); ip != "" {
			servicesByIP[ip] = append(servicesByIP[ip], services...)
		}
	}
}

func isTrafficInfraContainer(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case k8s.CaptureContainerName, k8s.NginxContainerName, k8s.FrpcContainerName:
		return true
	default:
		return false
	}
}

func uniqueSortedTrafficServices(services []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(services))
	for _, service := range services {
		service = strings.TrimSpace(service)
		if service == "" || seen[service] {
			continue
		}
		seen[service] = true
		result = append(result, service)
	}
	sort.Strings(result)
	return result
}

func dominantTrafficService(services []string) string {
	if len(services) == 0 {
		return ""
	}
	return services[0]
}

func trafficDirection(srcInternal, dstInternal bool) string {
	switch {
	case !srcInternal && dstInternal:
		return "ingress"
	case srcInternal && !dstInternal:
		return "egress"
	case srcInternal:
		return "internal"
	default:
		return "external"
	}
}

func trafficEdgeKind(srcInternal, dstInternal bool) string {
	switch {
	case srcInternal && dstInternal:
		return "internal"
	case srcInternal || dstInternal:
		return "boundary"
	default:
		return "external"
	}
}

func buildTrafficEdgeID(srcIP, dstIP, direction string) string {
	return fmt.Sprintf("%s>%s#%s", srcIP, dstIP, direction)
}

func buildTrafficNodeLabel(ip string, internal bool, services []string) string {
	if len(services) > 0 {
		return fmt.Sprintf("%s %s", dominantTrafficService(services), ip)
	}
	if internal {
		return fmt.Sprintf("Victim %s", ip)
	}
	if parsed, err := netip.ParseAddr(ip); err == nil {
		if parsed.IsLoopback() {
			return "Loopback"
		}
		if parsed.IsPrivate() {
			return fmt.Sprintf("Private %s", ip)
		}
		if parsed.IsMulticast() {
			return fmt.Sprintf("Multicast %s", ip)
		}
	}
	return ip
}

func trafficNodeKind(internal bool) string {
	if internal {
		return "victim"
	}
	return "peer"
}

func trafficNodeSide(ip string, srcInternal, dstInternal bool, internalIPs map[string]bool) string {
	if internalIPs[ip] {
		return "center"
	}
	if !srcInternal && dstInternal && !internalIPs[ip] {
		return "left"
	}
	if srcInternal && !dstInternal && !internalIPs[ip] {
		return "right"
	}
	return "orbit"
}

func trafficNodeZone(ip string, internalIPs map[string]bool) string {
	if internalIPs[ip] {
		return "victim"
	}
	parsed, err := netip.ParseAddr(ip)
	if err != nil {
		return "external"
	}
	switch {
	case parsed.IsPrivate():
		return "private"
	case parsed.IsLoopback():
		return "loopback"
	default:
		return "external"
	}
}

func normalizeTrafficProtocol(protocol string) string {
	protocol = strings.TrimSpace(protocol)
	if protocol == "" {
		return "Unknown"
	}
	return strings.ToUpper(protocol)
}

func normalizeTrafficSubtype(subtype string) string {
	subtype = strings.TrimSpace(subtype)
	if subtype == "" {
		return "Unknown"
	}
	subtype = strings.TrimPrefix(subtype, "LayerType")
	return strings.ToUpper(subtype)
}

func dominantTrafficKey(items map[string]int64) string {
	if len(items) == 0 {
		return ""
	}
	type kv struct {
		Key   string
		Value int64
	}
	ordered := make([]kv, 0, len(items))
	for key, value := range items {
		ordered = append(ordered, kv{Key: key, Value: value})
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Value != ordered[j].Value {
			return ordered[i].Value > ordered[j].Value
		}
		return ordered[i].Key < ordered[j].Key
	})
	return ordered[0].Key
}

func sortTrafficProtocolKeys(items map[string]int64) []string {
	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if items[keys[i]] != items[keys[j]] {
			return items[keys[i]] > items[keys[j]]
		}
		return keys[i] < keys[j]
	})
	return keys
}

func addTrafficProcess(items map[string]*trafficProcessAggregate, process *traffic.ProcessInfo, bytes int64) {
	if items == nil || process == nil {
		return
	}
	key := trafficProcessKey(*process)
	if key == "" {
		return
	}
	item := items[key]
	if item == nil {
		info := *process
		item = &trafficProcessAggregate{Info: info}
		items[key] = item
	}
	item.Bytes += bytes
	item.Packets++
}

func trafficProcessKey(process traffic.ProcessInfo) string {
	pid := ""
	if process.PID != nil {
		pid = strconv.FormatInt(*process.PID, 10)
	}
	name := strings.TrimSpace(process.ProcessName)
	if pid == "" && name == "" {
		return ""
	}
	return pid + ":" + name
}

func buildTrafficProcessResp(items map[string]*trafficProcessAggregate, limit int) []resp.TrafficProcessResp {
	if len(items) == 0 {
		return nil
	}
	processes := make([]*trafficProcessAggregate, 0, len(items))
	for _, item := range items {
		processes = append(processes, item)
	}
	sort.Slice(processes, func(i, j int) bool {
		if processes[i].Bytes != processes[j].Bytes {
			return processes[i].Bytes > processes[j].Bytes
		}
		return trafficProcessKey(processes[i].Info) < trafficProcessKey(processes[j].Info)
	})
	if limit <= 0 || len(processes) < limit {
		limit = len(processes)
	}
	result := make([]resp.TrafficProcessResp, 0, limit)
	for i := 0; i < limit; i++ {
		process := processes[i]
		result = append(result, resp.TrafficProcessResp{
			PID:              process.Info.PID,
			ProcessName:      process.Info.ProcessName,
			Bytes:            process.Bytes,
			Packets:          process.Packets,
			BytesSent:        process.Info.BytesSent,
			BytesReceived:    process.Info.BytesReceived,
			GeoIPCountryCode: process.Info.GeoIPCountryCode,
			GeoIPCountryName: process.Info.GeoIPCountryName,
			GeoIPASN:         process.Info.GeoIPASN,
			GeoIPASOrg:       process.Info.GeoIPASOrg,
			GeoIPCity:        process.Info.GeoIPCity,
			GeoIPPostalCode:  process.Info.GeoIPPostalCode,
		})
	}
	return result
}

func dominantTrafficProcess(processes []resp.TrafficProcessResp) string {
	if len(processes) == 0 {
		return ""
	}
	if processes[0].ProcessName != "" {
		return processes[0].ProcessName
	}
	if processes[0].PID != nil {
		return "PID " + strconv.FormatInt(*processes[0].PID, 10)
	}
	return ""
}

func countTrafficProcesses(connections []traffic.Connection) int {
	seen := make(map[string]bool)
	for _, connection := range connections {
		if connection.Process == nil {
			continue
		}
		key := trafficProcessKey(*connection.Process)
		if key != "" {
			seen[key] = true
		}
	}
	return len(seen)
}

func trafficIntensity(bytes, maxBytes int64) float64 {
	if maxBytes <= 0 {
		return 0
	}
	intensity := float64(bytes) / float64(maxBytes)
	if intensity < 0.15 {
		return 0.15
	}
	if intensity > 1 {
		return 1
	}
	return intensity
}

func computeTrafficPeak(buckets map[int64]*trafficBucketAggregate) (int64, int64) {
	peakTimeMs := int64(0)
	peakBytes := int64(0)
	for timestampMs, bucket := range buckets {
		if bucket.Bytes > peakBytes || (bucket.Bytes == peakBytes && timestampMs < peakTimeMs) {
			peakTimeMs = timestampMs
			peakBytes = bucket.Bytes
		}
	}
	return peakTimeMs, peakBytes
}

func availableTrafficSlices(totalDuration int64) []int64 {
	base := []int64{1000, 5000, 15000, 30000, 60000}
	if totalDuration <= 0 {
		return base
	}
	list := make([]int64, 0, len(base)+1)
	for _, candidate := range base {
		if candidate <= totalDuration {
			list = append(list, candidate)
		}
	}
	if len(list) == 0 {
		list = append(list, totalDuration)
	} else if list[len(list)-1] != totalDuration {
		list = append(list, totalDuration)
	}
	return list
}

func buildTrafficCenterLabel(victim model.Victim) string {
	if victim.ContestChallengeID.Valid && victim.ContestChallengeID.V > 0 {
		return "Victim #" + strconv.FormatUint(uint64(victim.ID), 10)
	}
	return "Instance #" + strconv.FormatUint(uint64(victim.ID), 10)
}

func buildTrafficRankingResp(items []trafficRankingAggregate, limit int) []resp.TrafficRankingResp {
	if limit <= 0 || len(items) < limit {
		limit = len(items)
	}
	result := make([]resp.TrafficRankingResp, 0, limit)
	for i := 0; i < limit; i++ {
		item := items[i]
		result = append(result, resp.TrafficRankingResp{
			Label:         item.Label,
			IP:            item.IP,
			Bytes:         item.Bytes,
			Packets:       item.Packets,
			Connections:   item.Connections,
			DominantProto: item.DominantProto,
			DominantApp:   item.DominantApp,
			DominantProc:  item.DominantProc,
			Direction:     item.Direction,
			Processes:     item.Processes,
		})
	}
	return result
}
