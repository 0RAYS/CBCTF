package k8s

import (
	"CBCTF/internal/model"
	"fmt"
)

func resourceError(ret model.RetVal) error {
	if detail, ok := ret.Attr["Error"]; ok {
		return fmt.Errorf("%s: %v", ret.Msg, detail)
	}
	return fmt.Errorf("%s", ret.Msg)
}
