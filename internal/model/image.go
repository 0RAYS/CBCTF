package model

import "time"

// NodeImageDetection is what one detection run found on one node.
// Present holds the images the node holds.
// Absent holds the images the node does not hold.
// Unverified holds the images the probe did not answer for.
type NodeImageDetection struct {
	Node       string   `json:"node"`
	Present    []string `json:"present"`
	Absent     []string `json:"absent"`
	Unverified []string `json:"unverified"`
}

// ImageDetection is the detection result of the common cache. Every run
// merges into it, so the system-wide view and the contest views read the same
// results.
type ImageDetection struct {
	DetectedAt time.Time            `json:"detected_at"`
	Nodes      []NodeImageDetection `json:"nodes"`
}
