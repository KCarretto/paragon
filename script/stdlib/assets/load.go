package assets

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kcarretto/paragon/script"
)

// Load uses assets.Open to load the asset and then it reads all the bytes.
//
// @param asset: A string for the path of the asset in the assets folder.
//
//
// @return (assetContent, nil) iff success; (nil, err) o/w
//
// @example
//  load("assets", "load")
//  load("sys", "write")
//
//  myBotAsset = load("/myBot/myBot.bin")
//  write("path/for/bot", myBotAsset)
func Load(assets http.FileSystem) script.Func {
	return script.Func(func(parser script.ArgParser) (script.Retval, error) {
		if assets == nil {
			return nil, fmt.Errorf("no assets available")
		}
		file, err := parser.GetString(0)
		if err != nil {
			return nil, err
		}
		filePtr, err := assets.Open(file)
		if err != nil {
			return nil, err
		}
		defer filePtr.Close()
		assetBin, err := ioutil.ReadAll(filePtr)
		if err != nil {
			return nil, err
		}
		return string(assetBin), nil
	})
}
