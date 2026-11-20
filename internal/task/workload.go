package task

import (
	"fmt"

	"CBCTF/internal/model"
)

// taskResourceError wraps a failed RetVal so that the resource kind, the
// target that was looked up and the underlying cause all reach the logs.
func taskResourceError(operation string, ret model.RetVal) error {
	return fmt.Errorf("%s: %w", operation, ret.AsError())
}
