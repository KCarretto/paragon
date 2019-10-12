package asset

import (
	"io/ioutil"

	"github.com/kcarretto/paragon/script"
)

// Load uses pkg.assets.Open to load the asset and then it reads all the bytes.
//
// @param asset: A string for the path of the asset in the asset folder.
//
//
// @return (assetBin, nil) iff success; (nil, err) o/w
//
// @example
//  load("assets", "load")
//  load("sys", "write")
//
//  myBotAsset = load("/myBot/myBot.bin")
//  write("path/for/bot", myBotAsset)
func Load(parser script.ArgParser) (script.Retval, error) {
	file, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	filePtr, err := pkg.assets.Open(file)
	if err != nil {
		return nil, err
	}
	defer filePtr.Close()
	assetBin, err := ioutil.ReadAll(filePtr)
	if err != nil {
		return nil, err
	}
	return string(assetBin), nil
}
