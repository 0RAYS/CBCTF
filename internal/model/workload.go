package model

// Workload status values shared by the dynamic challenge generators and the
// target instances, which follow the same lifecycle.
const (
	WaitingStatus     = "waiting"
	PendingStatus     = "pending"
	RunningStatus     = "running"
	TerminatingStatus = "terminating"
	StoppedStatus     = "stopped"
)

// IsStopping reports whether a workload in the given status is already
// terminating or has stopped.
func IsStopping(status string) bool {
	return status == TerminatingStatus || status == StoppedStatus
}
