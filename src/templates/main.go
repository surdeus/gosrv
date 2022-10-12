package templates

import(
	"html/template"
	"io/ioutil"
	"io"
	"log"
)

type Templates map[string] *template.Template
type FuncMap = template.FuncMap
type ParseConfig struct {
	Sep, Gen string
	FuncMap FuncMap
}

func (tmpls Templates)Execute(w io.Writer, t string, v any) {
	err := tmpls[t].ExecuteTemplate(w, t, v)
	if err != nil {
		log.Println(err)
	}
}

func MustParseTemplates(cfg ParseConfig) Templates {
	ret := make(map[string] *template.Template)

	files, err := ioutil.ReadDir(cfg.Sep)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		ret[f.Name()] =
			MustParseTemplate(cfg, f.Name())
	}

	return ret
}

func MustParseTemplate(cfg ParseConfig, name string) *template.Template {
	sep := cfg.Sep
	gen := cfg.Gen
	funcMap := cfg.FuncMap

	lfs := []string{sep + "/" + name}

	files, _ := ioutil.ReadDir(gen)
	for _, f := range files {
		lfs = append(lfs, gen+"/"+f.Name())
	}

	tmpl, err := template.New("").
		Funcs(funcMap).ParseFiles(lfs...)
	if err != nil {
		panic(err)
	}

	return tmpl
}

