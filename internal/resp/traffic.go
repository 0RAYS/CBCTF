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
	PeakSecond    int64 `json:"peak_second"`
	PeakBytes     int64 `json:"peak_bytes"`
}

type TrafficCenterResp struct {
	Label   string   `json:"label"`
	IPs     []string `json:"ips"`
	Exposed []string `json:"exposed"`
}

type TrafficNodeResp struct {
	ID            string   `json:"id"`
	Label         string   `json:"label"`
	IP            string   `json:"ip"`
	Kind          string   `json:"kind"`
	Side          string   `json:"side"`
	Zone          string   `json:"zone"`
	Bytes         int64    `json:"bytes"`
	Packets       int64    `json:"packets"`
	Connections   int64    `json:"connections"`
	Protocols     []string `json:"protocols"`
	DominantProto string   `json:"dominant_proto"`
}

type TrafficEdgeResp struct {
	ID            string   `json:"id"`
	Source        string   `json:"source"`
	Target        string   `json:"target"`
	Direction     string   `json:"direction"`
	Kind          string   `json:"kind"`
	Bytes         int64    `json:"bytes"`
	Packets       int64    `json:"packets"`
	Connections   int64    `json:"connections"`
	Weight        float64  `json:"weight"`
	Intensity     float64  `json:"intensity"`
	Protocols     []string `json:"protocols"`
	DominantProto string   `json:"dominant_proto"`
	DominantApp   string   `json:"dominant_app"`
}

type TrafficTimelineBucketResp struct {
	Second       int64 `json:"second"`
	Bytes        int64 `json:"bytes"`
	Packets      int64 `json:"packets"`
	IngressBytes int64 `json:"ingress_bytes"`
	EgressBytes  int64 `json:"egress_bytes"`
}

type TrafficRankingResp struct {
	Label         string `json:"label"`
	IP            string `json:"ip"`
	Bytes         int64  `json:"bytes"`
	Packets       int64  `json:"packets"`
	Connections   int64  `json:"connections"`
	DominantProto string `json:"dominant_proto"`
	DominantApp   string `json:"dominant_app,omitempty"`
	Direction     string `json:"direction,omitempty"`
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
}
