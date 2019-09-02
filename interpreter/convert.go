package interpreter

import "go.starlark.net/starlark"

func convertToStarlark(value interface{}) (starlark.Value, error) {
	if value == nil {
		return starlark.None, nil
	}
	switch v := value.(type) {
	case bool:
		return starlark.Bool(v), nil
	case int:
		return starlark.MakeInt(v), nil
	case int64:
		return starlark.MakeInt64(v), nil
	case uint:
		return starlark.MakeUint(v), nil
	case uint64:
		return starlark.MakeUint64(v), nil
	case float32:
		return starlark.Float(v), nil
	case float64:
		return starlark.Float(v), nil
	case string:
		return starlark.String(v), nil
	case []interface{}:
		var elems []starlark.Value
		for _, elem := range v {
			val, err := convertToStarlark(elem)
			if err != nil {
				return nil, err
			}
			elems = append(elems, val)
		}
		return starlark.NewList(elems), nil
	case map[interface{}]interface{}:
		dict := starlark.NewDict(len(v))

		for index, elem := range v {
			key, err := convertToStarlark(index)
			if err != nil {
				return nil, err
			}

			val, err := convertToStarlark(elem)
			if err != nil {
				return nil, err
			}

			if err = dict.SetKey(key, val); err != nil {
				return nil, err
			}
		}
		return dict, nil
	}
	return nil, nil
}
