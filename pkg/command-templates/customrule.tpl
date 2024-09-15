package customrule

import "github.com/olbrichattila/gofra/pkg/app/validator"

// {{.name}}Rule is a custom validator rule,
// val is the value to validate,
// pars is the elements in the rule signature, like myrule:1,2,3 will be 1, 2 and 3
// returns error message and bool if validation is OK
func {{.name}}Rule(val string, pars ...string) (validator.ValidationErrors, bool) {
	return nil, true
}
