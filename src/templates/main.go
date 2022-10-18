package templates

import(
	"html/template"
	"io"
	"log"
	"os"
	"strings"
	//"fmt"
	"path/filepath"
	"io/fs"
)

type Templates map[string] map[string] *template.Template

type FuncMap = template.FuncMap
type ParseConfig struct {
	View, Component, Template string
	FuncMap FuncMap
}

func (tmpls Templates)Exec(w io.Writer, t, v string, val any) {
	err := tmpls[t][v].Execute(w, val)
	if err != nil {
		log.Println(err)
	}
}

func Parse(cfg ParseConfig) (Templates, error) {
	var (
		ts Templates
		err error
	)
	t := template.New("").Funcs(cfg.FuncMap)

	t, err = parseFromDir(t, cfg.Component, "")
	if err != nil {
		return nil, err
	}
	

	ts = make(Templates)
	err = filepath.Walk(cfg.View,
		func(p string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() ||
				strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			viewName := p[len(cfg.View):]

			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}

			err = filepath.Walk(cfg.Template, func(
					pth string,
					info fs.FileInfo,
					err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() ||
					strings.HasPrefix(info.Name(), ".") {
					return nil
				}

				tv, err := t.Clone()
				if err != nil {
					return err
				}

				b2, err := os.ReadFile(pth)
				if err != nil {
					return err
				}

				tmplName := pth[len(cfg.Template):]
				//fmt.Printf("'%s' / '%s'\n", tmplName, viewName)

				_, ok := ts[tmplName]
				if !ok {
					ts[tmplName] = make(map[string] *template.Template)
				}

				st, err := tv.New("@master").Parse(string(b2))
				if err != nil {
					return err
				}
				st, err = st.New("").Parse(string(b))
				if err != nil {
					return err
				}

				ts[tmplName][viewName] = st


				return nil
			})

			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return ts, nil
}

func
parseFromDir(t *template.Template, root, dir string)(*template.Template, error) {
	p := root
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
			t, err = parseFromDir(t, root, tmplName)
			if err != nil {
				return nil, err
			}
		} else {
			if strings.HasPrefix(fileName, ".") {
				continue
			}

			b, err := os.ReadFile(filePath)
			if err != nil {
				return nil, err
			}

			//fmt.Printf("Parsing '%s'...\n", tmplName)
			t, err = t.New(tmplName).Parse(string(b))
			if err != nil {
					return nil, err
			}
		}
	}

	return t, nil
}

