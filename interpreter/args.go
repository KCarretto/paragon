package interpreter

import (
	"go.starlark.net/starlark"
)

// An ArgParser enables golang function implementations to retrieve the positional and keyword
// arguments provided to the starlark method call.
type ArgParser interface {
	GetString(index int) (string, error)
	GetStringByName(name string) (string, error)

	GetInt(index int) (int, error)
	GetIntByName(name string) (int, error)
}

type argParser struct {
	ArgParser

	args   map[int]interface{}
	kwargs map[string]int
}

func getParser(args starlark.Tuple, kwargs []starlark.Tuple) (*argParser, error) {
	parser := &argParser{
		args:   map[int]interface{}{},
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

func (parser *argParser) GetParam(index int) (interface{}, error) {
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

	str, ok := val.(string)
	if !ok {
		return "", ErrInvalidArgType
	}

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

	str, ok := val.(int)
	if !ok {
		return 0, ErrInvalidArgType
	}

	return str, nil
}
func (parser *argParser) GetIntByName(kwarg string) (int, error) {
	index, err := parser.GetParamIndex(kwarg)
	if err != nil {
		return 0, err
	}
	return parser.GetInt(index)
}
