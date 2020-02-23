// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

const SpecPath = "../../../../www/src/config/renegade/spec.json"

// ParamDefs enables easy parsing of parameter definition lists
type ParamDefs []string

// Parse the param definitions. Param defintions must be in the format "name@type".
func (paramDefs ParamDefs) Parse() (params []ParamSpec) {
	if paramDefs == nil {
		return
	}

	for _, paramDef := range paramDefs {
		parts := strings.Split(paramDef, "@")
		if len(parts) != 2 {
			panic(fmt.Errorf("provided invalid param spec %q, params should be in the form name@type",
				paramDef,
			))
		}
		param := ParamSpec{
			Name: parts[0],
			Type: parts[1],
		}
		params = append(params, param)
	}
	return
}

type ParamSpec struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type FunctionSpec struct {
	Name    string      `json:"name"`
	Doc     string      `json:"doc"`
	Params  []ParamSpec `json:"params"`
	Retvals []ParamSpec `json:"retvals"`
}

type LibrarySpec struct {
	Name      string          `json:"name"`
	Functions []*FunctionSpec `json:"functions"`
}

// GetFunction finds a function in the spec by name, and upserts it if it does not yet exist.
func (s *LibrarySpec) GetFunction(name string) *FunctionSpec {
	for _, fn := range s.Functions {
		if fn.Name == name {
			return fn
		}
	}

	fn := &FunctionSpec{
		Name:    name,
		Params:  []ParamSpec{},
		Retvals: []ParamSpec{},
	}

	s.Functions = append(s.Functions, fn)
	return fn
}

type Spec struct {
	Libraries []*LibrarySpec `json:"libraries"`
}

// GetLibrary finds a library in the spec by name, and upserts it if it does not yet exist.
func (s *Spec) GetLibrary(name string) *LibrarySpec {
	for _, lib := range s.Libraries {
		if lib.Name == name {
			return lib
		}
	}

	lib := &LibrarySpec{
		Name:      name,
		Functions: []*FunctionSpec{},
	}
	s.Libraries = append(s.Libraries, lib)
	return lib
}

func parseSpec(data []byte) (Spec, error) {
	var spec Spec
	if data == nil || len(data) < 1 {
		return spec, nil
	}

	if err := json.Unmarshal(data, &spec); err != nil {
		return Spec{}, fmt.Errorf("failed to parse spec file json: %w", err)
	}

	return spec, nil
}

func genDoc(specData []byte, lib, fn, doc string, params, retvals ParamDefs) (Spec, error) {
	spec, err := parseSpec(specData)
	if err != nil {
		return Spec{}, err
	}

	libSpec := spec.GetLibrary(lib)

	fnSpec := libSpec.GetFunction(fn)
	fnSpec.Doc = doc
	fnSpec.Params = params.Parse()
	fnSpec.Retvals = retvals.Parse()

	return spec, nil
}

func main() {
	params := &cli.StringSlice{}
	retvals := &cli.StringSlice{}

	var (
		lib string
		fn  string
		doc string
	)
	app := &cli.App{
		Name:  "gendoc",
		Usage: "Generate JSON documentation for languages",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "lib",
				Usage:       "Library that the function is defined in.",
				Destination: &lib,
			},
			&cli.StringFlag{
				Name:        "func",
				Value:       "func",
				Usage:       "Name of the function to document.",
				Destination: &fn,
			},
			&cli.StringFlag{
				Name:        "doc",
				Value:       "doc",
				Usage:       "Doc string to describe the function.",
				Destination: &doc,
			},
			&cli.StringSliceFlag{
				Name: "param",
				Usage: "Specify each param in the format of name@type. " +
					"If the function has no parameters, this option may be omitted.",
				Value: params,
			},
			&cli.StringSliceFlag{
				Name: "retval",
				Usage: "Specify each return value in the format of name@type. " +
					"If the function has no return values, this option may be omitted.",
				Value: retvals,
			},
		},
		Action: func(c *cli.Context) error {
			log.Printf("Generating spec: lib=%q func=%q doc=%q, params=%q, retvals=%q\n",
				lib,
				fn,
				doc,
				params.String(),
				retvals.String(),
			)
			// Get existing spec data, if any exists
			specData, err := ioutil.ReadFile(SpecPath)
			if err != nil {
				log.Printf("[WARN] Failed to read spec data, assuming no valid data available: %s", err.Error())
			}

			// Open and truncate the spec file, creating it if it doesn't yet exist
			f, err := os.OpenFile(SpecPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("failed to open spec file: %w", err)
			}
			defer f.Close()

			// Build the new spec
			spec, err := genDoc(specData, lib, fn, doc, ParamDefs(params.Value()), ParamDefs(retvals.Value()))

			// Marshal the new spec to JSON
			resultData, err := json.Marshal(spec)
			if err != nil {
				return fmt.Errorf("failed to marshal resulting spec json: %w", err)
			}

			// Write the spec to the file
			if _, err := f.Write(resultData); err != nil {
				return fmt.Errorf("failed to write new spec file: %w", err)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(fmt.Errorf("failed generation for %q: %w", fn, err))
	}
}
