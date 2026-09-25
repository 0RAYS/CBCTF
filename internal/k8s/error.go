package k8s

import (
	"fmt"

	"CBCTF/internal/model"
)

func resourceError(ret model.RetVal) error {
	if detail, ok := ret.Attr["Error"]; ok {
		return fmt.Errorf("%s: %v", ret.Msg, detail)
	}
	return fmt.Errorf("%s", ret.Msg)
}
