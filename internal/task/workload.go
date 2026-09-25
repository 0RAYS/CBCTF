package task

import (
	"fmt"

	"CBCTF/internal/model"
)

func taskResourceError(operation string, ret model.RetVal) error {
	if detail, ok := ret.Attr["Error"]; ok {
		return fmt.Errorf("%s: %s: %v", operation, ret.Msg, detail)
	}
	return fmt.Errorf("%s: %s", operation, ret.Msg)
}
