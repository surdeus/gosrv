package tmplx

import (
	"html/template"
	"fmt"
	"reflect"
)

type FuncConfig struct {
	StylePath string
	StyleFormat string
	ScriptPath string
	ScriptFormat string
}

func StdFuncMap() FuncMap {
	return FuncMap {
		"array" : Array,
		"string" : String,
		"hasField" : HasField,
		"attr" : Attr,
	}
}

func StdFuncCfg() FuncConfig {
	var (
		cfg FuncConfig
	)

	cfg.StylePath = "/s/style"
	cfg.ScriptPath = "/s/script"

	cfg.StyleFormat = "<link rel=\"stylesheet\" type=\"text/css\" href=\"%s/%s.css\">"
	cfg.ScriptFormat = "<script type=\"text/javascript\" src=\"%s/%s.js\"></script>"

	return cfg
}

func (cfg FuncConfig)Styles(styles ...string) template.HTML {
	ret := ""
	for _, style := range styles {
		ret += fmt.Sprintf(cfg.StyleFormat, cfg.StylePath, style)
	}

	return template.HTML(ret)
}

func (cfg FuncConfig)Scripts(scripts ...string) template.HTML {
	ret := ""
	for _, script := range scripts {
		ret += fmt.Sprintf(cfg.ScriptFormat, cfg.ScriptPath, script)
	}

	return template.HTML(ret)
}

func Array(args ...any) []any {
	return args
}

func Attr(args ...string) template.HTMLAttr {
	ret := ""
	for i, v := range args {
		ret += v
		if i != len(args)-1 {
			ret += " "
		}
	}

	return template.HTMLAttr(ret)
}

func String(v any) string {
	return v.(string)
}

func HasField(v any, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

func Sum(a, b int) int {
	return a + b
}

func Neg(a int) int {
	return -a
}

