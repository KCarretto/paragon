package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

// Assets contains project assets.
var Assets http.FileSystem = http.Dir("assets")

func main() {
	if err := vfsgen.Generate(Assets, vfsgen.Options{
		PackageName: "main",
		// BuildTags:    "",
		VariableName: "Assets",
		Filename:     "assets.gen.go",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := ioutil.WriteFile("main.gen.go", []byte(dropperMain), 0644); err != nil {
		log.Fatalln(err)
	}
}
