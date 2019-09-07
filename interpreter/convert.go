package interpreter

import (
	"reflect"

	"go.starlark.net/starlark"
)

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
	// case []interface{}:
	// 	var elems []starlark.Value
	// 	for _, elem := range v {
	// 		val, err := convertToStarlark(elem)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		elems = append(elems, val)
	// 	}
	// 	return starlark.NewList(elems), nil
	// case map[interface{}]interface{}:
	// 	dict := starlark.NewDict(len(v))

	// 	for index, elem := range v {
	// 		key, err := convertToStarlark(index)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		val, err := convertToStarlark(elem)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		if err = dict.SetKey(key, val); err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// 	return dict, nil
	default:
		reflectV := reflect.ValueOf(value)
		switch reflectV.Kind() {
		case reflect.Slice:
			var elems []starlark.Value
			for i := 0; i < reflectV.Len(); i++ {
				val, err := convertToStarlark(reflectV.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				elems = append(elems, val)
			}
			return starlark.NewList(elems), nil
		case reflect.Map:
			dict := starlark.NewDict(len(reflectV.MapKeys()))

			iter := reflectV.MapRange()
			for iter.Next() {
				key, err := convertToStarlark(iter.Key().Interface())
				if err != nil {
					return nil, err
				}

				val, err := convertToStarlark(iter.Value().Interface())
				if err != nil {
					return nil, err
				}
				if err = dict.SetKey(key, val); err != nil {
					return nil, err
				}
			}
			return dict, nil
		}
	}
	return nil, nil
}
