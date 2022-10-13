package templates

import(
	"html/template"
	"io"
	"log"
	"os"
)

type Templates struct {
	*template.Template
}
type FuncMap = template.FuncMap
type ParseConfig struct {
	Root string
	FuncMap FuncMap
}

func (tmpls *Templates)Exec(w io.Writer, t string, v any) {
	err := tmpls.ExecuteTemplate(w, t, v)
	if err != nil {
		log.Println(err)
	}
}

func Parse(cfg ParseConfig) (*Templates, error) {
	t := template.New("")

	t, err := parseFromDir(t, cfg, "")
	if err != nil {
		return nil, err
	}

	return &Templates{t}, nil
}

func
parseFromDir(t *template.Template, cfg ParseConfig, dir string)(*template.Template, error) {
	p := cfg.Root
	if dir != "" {
		p += "/" + dir
	}

	files, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileName := file.Name()
		filePath := fileName
		if filePath != "" {
			filePath = p + "/" + filePath
		}


		tmplName := fileName
		if dir != "" {
			tmplName = dir + "/" + tmplName
		}

		if file.IsDir() {
			t, err = parseFromDir(t, cfg, tmplName)
			if err != nil {
				return nil, err
			}
		} else {
			b, err := os.ReadFile(filePath)
			if err != nil {
				return nil, err
			}

			t = t.New(tmplName)
			t, err = t.Funcs(cfg.FuncMap).Parse(string(b))
			if err != nil {
					return nil, err
			}
		}
	}

	return t, nil
}

