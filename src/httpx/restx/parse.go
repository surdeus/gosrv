package restx

import (
	"strings"
	"net/url"
	//"fmt"
)

type ArgCfg struct {
	Sep string
}

type Args map[string] Arg

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
) Args {
	ret := make(Args)
	for k, v := range values {
		//fmt.Printf("%s %q, %q\n", k, v, strings.Split(k, ac.Sep))
		buf := Arg{}
		buf.Splits = strings.Split(k, ac.Sep)
		buf.Values = v
		ret[buf.Splits[0]] = buf
	}

	return ret
}

