package model

import "time"

// NodeImageDetection is what one detection run found on one node.
type NodeImageDetection struct {
	Node string `json:"node"`
	// Present holds the images the node holds.
	Present []string `json:"present"`
	// Absent holds the images the node does not hold.
	Absent []string `json:"absent"`
	// Unverified holds the images the probe did not answer for.
	Unverified []string `json:"unverified"`
}

// ImageDetection is the last detection result of one scope.
type ImageDetection struct {
	DetectedAt time.Time            `json:"detected_at"`
	Nodes      []NodeImageDetection `json:"nodes"`
}
