package interpreter

import (
	"go.starlark.net/starlark"
)

// TODO: Fix type assertions

// An ArgParser enables golang function implementations to retrieve the positional and keyword
// arguments provided to the starlark method call.
type ArgParser interface {
	GetString(index int) (string, error)
	GetStringByName(name string) (string, error)

	GetInt(index int) (int, error)
	GetIntByName(name string) (int, error)

	GetBool(index int) (bool, error)
	GetBoolByName(kwarg string) (bool, error)
}

type argParser struct {
	ArgParser

	args   map[int]starlark.Value
	kwargs map[string]int
}

func getParser(args starlark.Tuple, kwargs []starlark.Tuple) (*argParser, error) {
	parser := &argParser{
		args:   map[int]starlark.Value{},
		kwargs: map[string]int{},
	}

	for i, arg := range args {
		parser.args[i] = arg
	}
	for _, kwarg := range kwargs {
		if kwarg.Len() != 2 {
			return nil, ErrMalformattedKwarg
		}
		name := kwarg[0]
		val := kwarg[1]

		index := len(parser.args)
		parser.args[index] = val
		parser.kwargs[name.String()] = index
	}

	return parser, nil
}

func (parser *argParser) GetParam(index int) (starlark.Value, error) {
	val, ok := parser.args[index]
	if !ok {
		return nil, ErrMissingArg
	}
	return val, nil
}

func (parser *argParser) GetParamIndex(kwarg string) (int, error) {
	index, ok := parser.kwargs[kwarg]
	if !ok {
		return 0, ErrMissingKwarg
	}
	return index, nil
}

func (parser *argParser) GetString(index int) (string, error) {
	val, err := parser.GetParam(index)
	if err != nil {
		return "", err
	}

	if val.Type() != "string" {
		return "", ErrInvalidArgType
	}

	str, _ := starlark.AsString(val)
	return str, nil
}
func (parser *argParser) GetStringByName(kwarg string) (string, error) {
	index, err := parser.GetParamIndex(kwarg)
	if err != nil {
		return "", err
	}
	return parser.GetString(index)
}

func (parser *argParser) GetInt(index int) (int, error) {
	val, err := parser.GetParam(index)
	if err != nil {
		return 0, err
	}

	if val.Type() != "int" {
		return 0, ErrInvalidArgType
	}

	return starlark.AsInt32(val)
}
func (parser *argParser) GetIntByName(kwarg string) (int, error) {
	index, err := parser.GetParamIndex(kwarg)
	if err != nil {
		return 0, err
	}
	return parser.GetInt(index)
}

func (parser *argParser) GetBool(index int) (bool, error) {
	val, err := parser.GetParam(index)
	if err != nil {
		return false, err
	}

	if val.Type() != "bool" {
		return false, ErrInvalidArgType
	}

	return bool(val.Truth()), nil
}
func (parser *argParser) GetBoolByName(kwarg string) (bool, error) {
	index, err := parser.GetParamIndex(kwarg)
	if err != nil {
		return false, err
	}
	return parser.GetBool(index)
}
