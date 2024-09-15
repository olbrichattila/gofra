package internalcommand

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/olbrichattila/gofra/pkg/app/config"
)

func ListViewFunctions(c config.Configer) {
	viewFunctions := c.ViewConfig()

	keys := make([]string, 0, len(viewFunctions))
	for k := range viewFunctions {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, name := range keys {
		item := viewFunctions[name]
		fnSignature := getFunctionSignature(item)
		fmt.Printf("- %s [%s]\n", name, fnSignature)
	}
}

func getFunctionSignature(fn interface{}) string {
	funcType := reflect.TypeOf(fn)

	if funcType.Kind() != reflect.Func {
		return "Provided interface is not a function"
	}

	signature := "func("

	for i := 0; i < funcType.NumIn(); i++ {
		if i > 0 {
			signature += ", "
		}
		signature += funcType.In(i).String()
	}

	signature += ")"

	if funcType.NumOut() == 1 {
		signature += " " + funcType.Out(0).String()
	} else if funcType.NumOut() > 1 {
		signature += " ("
		for i := 0; i < funcType.NumOut(); i++ {
			if i > 0 {
				signature += ", "
			}
			signature += funcType.Out(i).String()
		}
		signature += ")"
	}

	return signature
}
