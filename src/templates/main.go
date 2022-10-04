package templates

import(
	"html/template"
	"io/ioutil"
	"io"
	"reflect"
	"log"
	//"fmt"
)

type Templates map[string] *template.Template
type FuncMap = template.FuncMap

func (tmpls Templates)Execute(w io.Writer, t string, v any) {
	err := tmpls[t].ExecuteTemplate(w, t, v)
	if err != nil {
		log.Println(err)
	}
}

func MustParseTemplates(sep, gen string, funcMap FuncMap) Templates {
	ret := make(map[string] *template.Template)

	files, err := ioutil.ReadDir(sep)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		ret[f.Name()] =
			MustParseTemplate(sep, gen, f.Name(), funcMap)
	}

	return ret
}

func MustParseTemplate(sep, gen, name string, funcMap template.FuncMap) *template.Template {
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

func HasField(v interface{}, name string) bool {
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

