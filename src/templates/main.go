package templates

import(
	"html/template"
	"io"
	"log"
	"os"
	"strings"
	"fmt"
	"path/filepath"
	"io/fs"
)

type Templates map[string] *template.Template

type FuncMap = template.FuncMap
type ParseConfig struct {
	Gen, Sep string
	FuncMap FuncMap
}

func (tmpls Templates)Exec(w io.Writer, t string, v any) {
	err := tmpls[t].Execute(w, v)
	if err != nil {
		log.Println(err)
	}
}

func Parse(cfg ParseConfig) (Templates, error) {
	var (
		ts Templates
	)
	t := template.New("").Funcs(cfg.FuncMap)

	t, err := parseFromDir(t, cfg.Gen, "")
	if err != nil {
		return nil, err
	}

	ts = make(Templates)
	err = filepath.Walk(cfg.Sep,
		func(p string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() ||
				strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			tmplName := p[len(cfg.Sep):]
			fmt.Printf("'%s'\n", tmplName)

			tv, err := t.Clone()
			if err != nil {
				return err
			}

			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}

			ts[tmplName], err = tv.
				New("").
				Parse(string(b))
			return nil
		})

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

