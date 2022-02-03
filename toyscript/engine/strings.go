package engine

import (
	"fmt"
	"strconv"
)

func (value *ToyScriptValue) ToString() string {
	switch value.Type {
	case TString:
		return value.Value.(string)
	case TInt:
		// TODO - move away from builtin
		return strconv.FormatInt(value.Value.(int64), 10)
	case TBool:
		if value.Value == true {
			return "true"
		} else {
			return "false"
		}
	case TFloat:
		return strconv.FormatFloat(value.Value.(float64), 'f', -1, 64)
	case TNull:
		return "<null>"
	}
	return fmt.Sprint(value)
}
