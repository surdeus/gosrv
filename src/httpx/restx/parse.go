package restx

import (
	"strings"
	"net/url"
	//"fmt"
)

type ArgCfg struct {
	Sep string
}

type Arg struct {
	Splits []string
	Values []string
}

func StdArgCfg() *ArgCfg {
	ret := ArgCfg {
		Sep : "__",
	}
	return &ret
}

func (ac *ArgCfg) ParseValues (
	values url.Values,
) []Arg {
	ret := []Arg{}
	for k, v := range values {
		//fmt.Printf("%s %q, %q\n", k, v, strings.Split(k, ac.Sep))
		buf := Arg{}
		buf.Splits = strings.Split(k, ac.Sep)
		buf.Values = v
		ret = append(ret, buf)
	}

	return ret
}
