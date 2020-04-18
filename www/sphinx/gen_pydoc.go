// + build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

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
type Spec struct {
	Libraries []*LibrarySpec `json:"libraries"`
}

const SpecPath = "./src/config/renegade/spec.json"
const SphinxDir = "./sphinx/"
const RSTDir = SphinxDir + "source/"
const PublicDir = "./public/"

func makeParamsRST(params []ParamSpec) string {
	var paramStrings []string
	for _, p := range params {
		paramStrings = append(paramStrings, p.Name+": "+p.Type)
	}
	return strings.Join(paramStrings, ",")
}

func makeFuncsRST(funcs []*FunctionSpec) string {
	ret := ""
	for _, fun := range funcs {
		sig := ".. function:: " + fun.Name + "(" + makeParamsRST(fun.Params) + ") -> (" + makeParamsRST(fun.Retvals) + ")\n\n"
		headerTwoDel := "----\n\n"
		ret = ret + sig + "\t" + fun.Doc + "\n\n" + headerTwoDel
	}
	return ret
}

func makeLibsRST(libs []*LibrarySpec) string {
	ret := ""
	for _, lib := range libs {
		libHeader := "stdlib/" + lib.Name + "\n"
		headerOneDel := "--------------------------------------\n\n"
		currModule := ".. currentmodule:: " + lib.Name + "\n\n"
		ret = ret + libHeader + headerOneDel + currModule + makeFuncsRST(lib.Functions)
	}
	return ret
}

func main() {
	specData, err := ioutil.ReadFile(SpecPath)
	if err != nil {
		panic("[ERR] Failed to read spec data, assuming no valid data available")
	}
	var spec Spec
	err = json.Unmarshal(specData, &spec)
	if err != nil {
		panic("[ERR] Failed to parse spec data, assuming no valid data available")
	}

	indexRST := `
Welcome to Renegade's documentation!
====================================
The following doc page is automatically generated. Please be kind to it.


`

	indexRST = indexRST + makeLibsRST(spec.Libraries)
	err = ioutil.WriteFile(RSTDir+"index.rst", []byte(indexRST), 0644)
	if err != nil {
		panic("Failed to write out RST File paniccccccccc")
	}

	err = os.Chdir(SphinxDir)
	if err != nil {
		panic("Failed to change directories for sphinx")
	}

	cmd := exec.Command("/usr/bin/make", "html")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic("Failed to run make for sphinx: " + string(out))
	}
	fmt.Printf("%s\n", out)

	os.RemoveAll("../" + PublicDir + "docs/")

	err = os.Rename("build/html/", "../"+PublicDir+"docs/")
	if err != nil {
		panic("Failed to move docs to public folder: " + err.Error())
	}

	fmt.Printf("Sphinx Renegade Docs written\n")

}
