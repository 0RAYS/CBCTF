package task

import (
	"CBCTF/internal/model"
	"fmt"
)

func taskResourceError(operation string, ret model.RetVal) error {
	if detail, ok := ret.Attr["Error"]; ok {
		return fmt.Errorf("%s: %s: %v", operation, ret.Msg, detail)
	}
	return fmt.Errorf("%s: %s", operation, ret.Msg)
}
