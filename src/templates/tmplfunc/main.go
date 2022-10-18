package tmplfunc

import (
	"html/template"
	"fmt"
)

type Config struct {
	StylePath string
	StyleFormat string
	ScriptPath string
	ScriptFormat string
}

func StdCfg() Config {
	var (
		cfg Config
	)

	cfg.StylePath = "/s/style"
	cfg.ScriptPath = "/s/script"

	cfg.StyleFormat = "<link rel=\"stylesheet\" type=\"text/css\" href=\"%s/%s.css\">"
	cfg.ScriptFormat = "<script type=\"text/javascript\" src=\"%s/%s.js\"></script>"

	return cfg
}

func (cfg Config)Styles(styles ...string) template.HTML {
	ret := ""
	for _, style := range styles {
		ret += fmt.Sprintf(cfg.StyleFormat, cfg.StylePath, style)
	}

	return template.HTML(ret)
}

func (cfg Config)Scripts(scripts ...string) template.HTML {
	ret := ""
	for _, script := range scripts {
		ret += fmt.Sprintf(cfg.ScriptFormat, cfg.ScriptPath, script)
	}

	return template.HTML(ret)
}
