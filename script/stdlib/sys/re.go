package sys

import (
	"regexp"

	"github.com/kcarretto/paragon/script"
)

// ReplaceString uses regexp.MustCompile to replace values in a string.
//
// @param inStr:   A string to have parts replaced.
//
// @param newStr:  A string that will replace the places where the pattern match in `inStr`. Also can reference backrefs
// with `${1}`-like syntax.
//
// @param pattern: A string for the regular expression to pattern match with/generate capture groups.
//
//
// @return (resultStr, nil) iff success; (nil, err) o/w
//
// @example
//  load("sys", "read")
//  load("sys", "write")
//  load("sys", "replaceString")
//  load("sys", "exec")
//
//  nginx_conf_file = "/etc/nginx.conf"
//  nginx_conf = read(nginx_conf_file)
//  new_nginx_conf = replaceString(nginx_conf, "listen[\s]*80;", "listen 81;")
//  write(nginx_conf_file, new_nginx_conf)
// 	exec("systemctl restart nginx")
func ReplaceString(parser script.ArgParser) (script.Retval, error) {
	inStr, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	newStr, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	pattern, err := parser.GetString(2)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(inStr, newStr), nil
}
