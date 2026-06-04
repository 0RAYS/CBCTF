package resp

type TrafficWindowResp struct {
	Start      int64 `json:"start"`
	End        int64 `json:"end"`
	Duration   int64 `json:"duration"`
	Total      int64 `json:"total"`
	TotalCount int64 `json:"total_count"`
}

type TrafficSummaryResp struct {
	TotalBytes    int64 `json:"total_bytes"`
	TotalPackets  int64 `json:"total_packets"`
	IngressBytes  int64 `json:"ingress_bytes"`
	EgressBytes   int64 `json:"egress_bytes"`
	InternalBytes int64 `json:"internal_bytes"`
	ExternalNodes int   `json:"external_nodes"`
	InternalNodes int   `json:"internal_nodes"`
	VisibleEdges  int   `json:"visible_edges"`
	VisibleNodes  int   `json:"visible_nodes"`
	PeakTimeMs    int64 `json:"peak_time_ms"`
	PeakBytes     int64 `json:"peak_bytes"`
	ProcessCount  int   `json:"process_count"`
}

type TrafficProcessResp struct {
	PID              *int64 `json:"pid,omitempty"`
	ProcessName      string `json:"process_name,omitempty"`
	Bytes            int64  `json:"bytes"`
	Packets          int64  `json:"packets"`
	BytesSent        int64  `json:"bytes_sent,omitempty"`
	BytesReceived    int64  `json:"bytes_received,omitempty"`
	GeoIPCountryCode string `json:"geoip_country_code,omitempty"`
	GeoIPCountryName string `json:"geoip_country_name,omitempty"`
	GeoIPASN         *int64 `json:"geoip_asn,omitempty"`
	GeoIPASOrg       string `json:"geoip_as_org,omitempty"`
	GeoIPCity        string `json:"geoip_city,omitempty"`
	GeoIPPostalCode  string `json:"geoip_postal_code,omitempty"`
}

type TrafficCenterResp struct {
	Label   string   `json:"label"`
	Exposed []string `json:"exposed"`
}

type TrafficNodeResp struct {
	ID            string               `json:"id"`
	Label         string               `json:"label"`
	IP            string               `json:"ip"`
	Kind          string               `json:"kind"`
	Side          string               `json:"side"`
	Zone          string               `json:"zone"`
	Service       string               `json:"service,omitempty"`
	Services      []string             `json:"services,omitempty"`
	Bytes         int64                `json:"bytes"`
	Packets       int64                `json:"packets"`
	Connections   int64                `json:"connections"`
	Protocols     []string             `json:"protocols"`
	DominantProto string               `json:"dominant_proto"`
	DominantProc  string               `json:"dominant_process,omitempty"`
	Processes     []TrafficProcessResp `json:"processes,omitempty"`
}

type TrafficEdgeResp struct {
	ID            string               `json:"id"`
	Source        string               `json:"source"`
	Target        string               `json:"target"`
	Direction     string               `json:"direction"`
	Kind          string               `json:"kind"`
	SourceService string               `json:"source_service,omitempty"`
	TargetService string               `json:"target_service,omitempty"`
	Bytes         int64                `json:"bytes"`
	Packets       int64                `json:"packets"`
	Connections   int64                `json:"connections"`
	Weight        float64              `json:"weight"`
	Intensity     float64              `json:"intensity"`
	Protocols     []string             `json:"protocols"`
	DominantProto string               `json:"dominant_proto"`
	DominantApp   string               `json:"dominant_app"`
	DominantProc  string               `json:"dominant_process,omitempty"`
	Processes     []TrafficProcessResp `json:"processes,omitempty"`
}

type TrafficTimelineBucketResp struct {
	TimestampMs  int64 `json:"timestamp_ms"`
	Bytes        int64 `json:"bytes"`
	Packets      int64 `json:"packets"`
	IngressBytes int64 `json:"ingress_bytes"`
	EgressBytes  int64 `json:"egress_bytes"`
}

type TrafficRankingResp struct {
	Label         string               `json:"label"`
	IP            string               `json:"ip"`
	Bytes         int64                `json:"bytes"`
	Packets       int64                `json:"packets"`
	Connections   int64                `json:"connections"`
	DominantProto string               `json:"dominant_proto"`
	DominantApp   string               `json:"dominant_app,omitempty"`
	DominantProc  string               `json:"dominant_process,omitempty"`
	Direction     string               `json:"direction,omitempty"`
	Processes     []TrafficProcessResp `json:"processes,omitempty"`
}

type TrafficTopologyResp struct {
	Window          TrafficWindowResp           `json:"window"`
	TotalDuration   int64                       `json:"total_duration"`
	AvailableSlices []int64                     `json:"available_slices"`
	Center          TrafficCenterResp           `json:"center"`
	Summary         TrafficSummaryResp          `json:"summary"`
	Nodes           []TrafficNodeResp           `json:"nodes"`
	Edges           []TrafficEdgeResp           `json:"edges"`
	Timeline        []TrafficTimelineBucketResp `json:"timeline"`
	TopTalkers      []TrafficRankingResp        `json:"top_talkers"`
	TopEdges        []TrafficRankingResp        `json:"top_edges"`
	IPs             []string                    `json:"ips"`
}
