package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	r "CBCTF/internal/redis"
	"CBCTF/internal/resp"
	"CBCTF/internal/utils"
	"fmt"
	"net/netip"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

type trafficNodeAggregate struct {
	ID            string
	Label         string
	IP            string
	Kind          string
	Side          string
	Zone          string
	Bytes         int64
	Packets       int64
	Connections   int64
	protocolBytes map[string]int64
}

type trafficEdgeAggregate struct {
	ID            string
	Source        string
	Target        string
	Direction     string
	Kind          string
	Bytes         int64
	Packets       int64
	Connections   int64
	protocolBytes map[string]int64
	appBytes      map[string]int64
}

type trafficBucketAggregate struct {
	Second       int64
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
	Direction     string
}

func GetTraffic(victim model.Victim, form dto.GetTrafficForm) (resp.TrafficTopologyResp, model.RetVal) {
	connections, ret := loadTrafficConnections(victim)
	if !ret.OK {
		return resp.TrafficTopologyResp{}, ret
	}
	if len(connections) == 0 {
		return emptyTrafficTopology(victim, form), model.SuccessRetVal()
	}

	totalDuration := calcTrafficTotalDuration(connections)
	start, end := clampTrafficWindow(form.TimeShift, form.Duration, totalDuration)
	windowConnections := sliceTrafficConnections(connections, start, end)

	internalIPs := collectVictimIPs(victim, connections)
	totalPackets := int64(len(windowConnections))

	nodes := make(map[string]*trafficNodeAggregate)
	edges := make(map[string]*trafficEdgeAggregate)
	buckets := make(map[int64]*trafficBucketAggregate)

	summary := resp.TrafficSummaryResp{}
	topTalkers := make([]trafficRankingAggregate, 0)
	topEdges := make([]trafficRankingAggregate, 0)

	for _, connection := range windowConnections {
		second := int64(connection.TimeShift / time.Second)
		bucket := buckets[second]
		if bucket == nil {
			bucket = &trafficBucketAggregate{Second: second}
			buckets[second] = bucket
		}
		packetBytes := int64(connection.Size)
		bucket.Bytes += packetBytes
		bucket.Packets++
		summary.TotalBytes += packetBytes
		summary.TotalPackets++

		srcInternal := internalIPs[connection.SrcIP]
		dstInternal := internalIPs[connection.DstIP]

		direction := trafficDirection(srcInternal, dstInternal)
		if direction == "ingress" {
			bucket.IngressBytes += packetBytes
			summary.IngressBytes += packetBytes
		} else if direction == "egress" {
			bucket.EgressBytes += packetBytes
			summary.EgressBytes += packetBytes
		} else {
			summary.InternalBytes += packetBytes
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
				protocolBytes: make(map[string]int64),
				appBytes:      make(map[string]int64),
			}
			edges[edgeID] = edge
		}
		edge.Bytes += packetBytes
		edge.Packets++
		edge.Connections++
		edge.protocolBytes[protocol] += packetBytes
		edge.appBytes[app] += packetBytes

		for _, ip := range []string{connection.SrcIP, connection.DstIP} {
			node := nodes[ip]
			if node == nil {
				node = &trafficNodeAggregate{
					ID:            ip,
					Label:         buildTrafficNodeLabel(ip, internalIPs[ip]),
					IP:            ip,
					Kind:          trafficNodeKind(internalIPs[ip]),
					Side:          trafficNodeSide(ip, srcInternal, dstInternal, internalIPs),
					Zone:          trafficNodeZone(ip, internalIPs),
					protocolBytes: make(map[string]int64),
				}
				nodes[ip] = node
			}
			node.Bytes += packetBytes
			node.Packets++
			node.Connections++
			node.protocolBytes[protocol] += packetBytes
		}
	}

	summary.PeakSecond, summary.PeakBytes = computeTrafficPeak(buckets)

	nodeList := make([]resp.TrafficNodeResp, 0, len(nodes))
	internalNodeCount := 0
	externalNodeCount := 0
	for _, node := range nodes {
		protocols := sortTrafficProtocolKeys(node.protocolBytes)
		nodeList = append(nodeList, resp.TrafficNodeResp{
			ID:            node.ID,
			Label:         node.Label,
			IP:            node.IP,
			Kind:          node.Kind,
			Side:          node.Side,
			Zone:          node.Zone,
			Bytes:         node.Bytes,
			Packets:       node.Packets,
			Connections:   node.Connections,
			Protocols:     protocols,
			DominantProto: dominantTrafficKey(node.protocolBytes),
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
		edgeList = append(edgeList, resp.TrafficEdgeResp{
			ID:            edge.ID,
			Source:        edge.Source,
			Target:        edge.Target,
			Direction:     edge.Direction,
			Kind:          edge.Kind,
			Bytes:         edge.Bytes,
			Packets:       edge.Packets,
			Connections:   edge.Connections,
			Weight:        float64(edge.Bytes) / float64(maxEdgeBytes),
			Intensity:     trafficIntensity(edge.Bytes, maxEdgeBytes),
			Protocols:     protocols,
			DominantProto: dominantProto,
			DominantApp:   dominantApp,
		})
		topEdges = append(topEdges, trafficRankingAggregate{
			Label:         fmt.Sprintf("%s -> %s", edge.Source, edge.Target),
			IP:            edge.ID,
			Bytes:         edge.Bytes,
			Packets:       edge.Packets,
			Connections:   edge.Connections,
			DominantProto: dominantProto,
			DominantApp:   dominantApp,
			Direction:     edge.Direction,
		})
	}
	sort.Slice(edgeList, func(i, j int) bool {
		if edgeList[i].Bytes != edgeList[j].Bytes {
			return edgeList[i].Bytes > edgeList[j].Bytes
		}
		return edgeList[i].ID < edgeList[j].ID
	})

	timeline := make([]resp.TrafficTimelineBucketResp, 0, len(buckets))
	seconds := make([]int64, 0, len(buckets))
	for second := range buckets {
		seconds = append(seconds, second)
	}
	sort.Slice(seconds, func(i, j int) bool { return seconds[i] < seconds[j] })
	for _, second := range seconds {
		bucket := buckets[second]
		timeline = append(timeline, resp.TrafficTimelineBucketResp{
			Second:       bucket.Second,
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

	return resp.TrafficTopologyResp{
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
			IPs:     sortedTrafficIPs(internalIPs),
			Exposed: victim.RemoteAddr(),
		},
		Summary:    summary,
		Nodes:      nodeList,
		Edges:      edgeList,
		Timeline:   timeline,
		TopTalkers: buildTrafficRankingResp(topTalkers, 6),
		TopEdges:   buildTrafficRankingResp(topEdges, 6),
	}, model.SuccessRetVal()
}

func loadTrafficConnections(victim model.Victim) ([]utils.Connection, model.RetVal) {
	connections, ret := r.GetTraffic(victim)
	if !ret.OK {
		return nil, ret
	}
	if len(connections) > 0 {
		return connections, model.SuccessRetVal()
	}
	ret = r.UpdateTraffics(victim)
	if !ret.OK {
		return nil, ret
	}
	connections, ret = r.GetTraffic(victim)
	if !ret.OK {
		return nil, ret
	}
	if len(connections) == 0 {
		return make([]utils.Connection, 0), model.SuccessRetVal()
	}
	return connections, model.SuccessRetVal()
}

func emptyTrafficTopology(victim model.Victim, form dto.GetTrafficForm) resp.TrafficTopologyResp {
	duration := form.Duration
	if duration <= 0 {
		duration = 15
	}
	return resp.TrafficTopologyResp{
		Window: resp.TrafficWindowResp{
			Start:    form.TimeShift,
			End:      form.TimeShift + duration,
			Duration: duration,
			Total:    0,
		},
		TotalDuration:   0,
		AvailableSlices: []int64{5, 15, 30, 60},
		Center: resp.TrafficCenterResp{
			Label:   buildTrafficCenterLabel(victim),
			IPs:     sortedTrafficIPs(collectVictimIPs(victim, nil)),
			Exposed: victim.RemoteAddr(),
		},
		Summary:    resp.TrafficSummaryResp{},
		Nodes:      make([]resp.TrafficNodeResp, 0),
		Edges:      make([]resp.TrafficEdgeResp, 0),
		Timeline:   make([]resp.TrafficTimelineBucketResp, 0),
		TopTalkers: make([]resp.TrafficRankingResp, 0),
		TopEdges:   make([]resp.TrafficRankingResp, 0),
	}
}

func calcTrafficTotalDuration(connections []utils.Connection) int64 {
	if len(connections) == 0 {
		return 0
	}
	duration := int64(connections[len(connections)-1].Time.Sub(connections[0].Time) / time.Second)
	return duration + 1
}

func clampTrafficWindow(start, duration, total int64) (int64, int64) {
	if start < 0 {
		start = 0
	}
	if duration <= 0 {
		duration = 15
	}
	if total > 0 && start > total {
		start = total
	}
	end := start + duration
	if total > 0 && end > total {
		end = total
	}
	if end < start {
		end = start
	}
	return start, end
}

func sliceTrafficConnections(connections []utils.Connection, start, end int64) []utils.Connection {
	if len(connections) == 0 {
		return make([]utils.Connection, 0)
	}
	startAt := time.Duration(start) * time.Second
	endAt := time.Duration(end) * time.Second
	if end == start {
		endAt = startAt + time.Second
	}

	windowConnections := make([]utils.Connection, 0)
	for _, connection := range connections {
		if connection.TimeShift < startAt {
			continue
		}
		if connection.TimeShift > endAt {
			break
		}
		windowConnections = append(windowConnections, connection)
	}
	return windowConnections
}

func collectVictimIPs(victim model.Victim, connections []utils.Connection) map[string]bool {
	internalIPs := make(map[string]bool)
	for _, pod := range victim.Pods {
		for _, network := range pod.Spec.Networks {
			if network.IP != "" {
				internalIPs[network.IP] = true
			}
		}
	}
	for _, pod := range victim.Spec.Pods {
		for _, network := range pod.Networks {
			if network.IP != "" {
				internalIPs[network.IP] = true
			}
		}
	}
	for _, endpoint := range victim.Endpoints {
		if endpoint.IP != "" {
			internalIPs[endpoint.IP] = true
		}
	}
	for _, endpoint := range victim.ExposedEndpoints {
		if endpoint.IP != "" {
			internalIPs[endpoint.IP] = true
		}
	}
	if len(internalIPs) > 0 {
		return internalIPs
	}

	freq := make(map[string]int)
	for _, connection := range connections {
		freq[connection.SrcIP]++
		freq[connection.DstIP]++
	}
	type kv struct {
		IP    string
		Count int
	}
	ordered := make([]kv, 0, len(freq))
	for ip, count := range freq {
		ordered = append(ordered, kv{IP: ip, Count: count})
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Count != ordered[j].Count {
			return ordered[i].Count > ordered[j].Count
		}
		return ordered[i].IP < ordered[j].IP
	})
	for _, item := range ordered {
		internalIPs[item.IP] = true
		if len(internalIPs) >= 3 {
			break
		}
	}
	return internalIPs
}

func trafficDirection(srcInternal, dstInternal bool) string {
	switch {
	case !srcInternal && dstInternal:
		return "ingress"
	case srcInternal && !dstInternal:
		return "egress"
	case srcInternal && dstInternal:
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

func buildTrafficNodeLabel(ip string, internal bool) string {
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
	peakSecond := int64(0)
	peakBytes := int64(0)
	for second, bucket := range buckets {
		if bucket.Bytes > peakBytes || (bucket.Bytes == peakBytes && second < peakSecond) {
			peakSecond = second
			peakBytes = bucket.Bytes
		}
	}
	return peakSecond, peakBytes
}

func availableTrafficSlices(totalDuration int64) []int64 {
	base := []int64{5, 15, 30, 60}
	if totalDuration <= 0 {
		return base
	}
	slices := make([]int64, 0, len(base)+1)
	for _, candidate := range base {
		if candidate <= totalDuration {
			slices = append(slices, candidate)
		}
	}
	if len(slices) == 0 {
		slices = append(slices, totalDuration)
	} else if slices[len(slices)-1] != totalDuration {
		slices = append(slices, totalDuration)
	}
	return slices
}

func buildTrafficCenterLabel(victim model.Victim) string {
	if victim.ContestChallengeID.Valid && victim.ContestChallengeID.V > 0 {
		return fmt.Sprintf("Victim #%d", victim.ID)
	}
	return fmt.Sprintf("Instance #%d", victim.ID)
}

func sortedTrafficIPs(items map[string]bool) []string {
	ips := make([]string, 0, len(items))
	for ip := range items {
		ips = append(ips, ip)
	}
	sort.Strings(ips)
	return ips
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
			Direction:     item.Direction,
		})
	}
	return result
}

// LoadTraffic 简单记录涉及到的 IP 地址
func LoadTraffic(tx *gorm.DB, victim model.Victim) model.RetVal {
	trafficRepo := db.InitTrafficRepo(tx)
	optionsL := make(map[string]db.CreateTrafficOptions)
	count, _ := trafficRepo.Count(db.CountOptions{Conditions: map[string]any{"victim_id": victim.ID}})
	if count > 0 {
		return model.SuccessRetVal()
	}
	go func(victim model.Victim) {
		if err := utils.Zip(victim.TrafficBasePath(), victim.TrafficZipPath()); err != nil {
			log.Logger.Warningf("Failed to zip .pcap files: %s", err)
			return
		}
		size, hash, err := utils.GetFileInfoByPath(victim.TrafficZipPath())
		if err != nil {
			log.Logger.Warningf("Failed to get file info %s: %s", victim.TrafficZipPath(), err)
			return
		}
		db.InitFileRepo(db.DB).Create(db.CreateFileOptions{
			RandID:   utils.UUID(),
			Filename: "traffics.zip",
			Size:     size,
			Path:     model.FilePath(victim.TrafficZipPath()),
			Model:    model.ModelName(victim),
			ModelID:  victim.ID,
			Suffix:   ".zip",
			Hash:     hash,
			Type:     model.TrafficFileType,
		})
	}(victim)
	connections, err := utils.ReadPcapDir(victim.TrafficBasePath())
	if err != nil {
		log.Logger.Warningf("Failed to read pcap: %s", err)
		return model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}
	for _, conn := range connections {
		connID := fmt.Sprintf("%s-%s-%s-%s", conn.SrcIP, conn.DstIP, conn.Type, conn.Subtype)
		if options, ok := optionsL[connID]; ok {
			options.Count += 1
			options.Size += conn.Size
			optionsL[connID] = options
		} else {
			optionsL[connID] = db.CreateTrafficOptions{
				VictimID: victim.ID,
				SrcIP:    conn.SrcIP,
				DstIP:    conn.DstIP,
				Type:     conn.Type,
				Subtype:  conn.Subtype,
				Size:     conn.Size,
				Count:    1,
			}
		}
	}
	for _, options := range optionsL {
		_, ret := trafficRepo.Create(options)
		if !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}
