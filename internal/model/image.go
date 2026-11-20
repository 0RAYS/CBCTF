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

// ImageDetection is the last detection result of one scope.
type ImageDetection struct {
	DetectedAt time.Time            `json:"detected_at"`
	Nodes      []NodeImageDetection `json:"nodes"`
}
