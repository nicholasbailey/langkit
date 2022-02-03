package engine

func BoolFromGoBoolean(x bool) *ToyScriptValue {
	if x {
		return True()
	}
	return False()
}

func False() *ToyScriptValue {
	return &ToyScriptValue{
		Type:  TBool,
		Value: false,
	}
}

func True() *ToyScriptValue {
	return &ToyScriptValue{
		Type:  TBool,
		Value: true,
	}
}

func Truthiness(value *ToyScriptValue) *ToyScriptValue {
	switch value.Type {
	case TBool:
		return value
	case TString:
		if value.Value.(string) == "" {
			return False()
		} else {
			return True()
		}
	case TInt:
		if value.Value.(int64) == 0 {
			return False()
		} else {
			return True()
		}
	case TFloat:
		if value.Value.(float64) == 0.0 {
			return False()
		} else {
			return True()
		}
	case TNull:
		return False()
	}
	panic("How did we get here")
}
