package task

import (
	"errors"
	"fmt"
	"testing"

	"CBCTF/internal/k8s"
)

func TestReadinessReadOutagesDoNotTriggerDestructiveCleanup(t *testing.T) {
	transient := errors.New("API/cache read unavailable")
	terminal := fmt.Errorf("pod state: %w", &k8s.StartupFailure{Err: errors.New("container exited")})
	for _, tc := range []struct {
		expired bool
		err     error
		stop    bool
	}{
		{false, nil, false}, {false, transient, false}, {false, terminal, true}, {true, transient, true}, {true, nil, true},
	} {
		if got := shouldStopUnreadyVictim(tc.expired, tc.err); got != tc.stop {
			t.Fatalf("expired=%t err=%v cleanup=%t", tc.expired, tc.err, got)
		}
	}
}
