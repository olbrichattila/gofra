package controller
{{.imports}}{{.structDef}}
// Index{{.name}}Action function can take any parameters defined in the Di config
func {{.receiver}}Index{{.name}}Action({{.in}}){{.out}} {
{{.return}}}

// Store{{.name}}Action function can take any parameters defined in the Di config
func {{.receiver}}Store{{.name}}Action({{.in}}){{.out}} {
{{.return}}}

// Show{{.name}}Action function can take any parameters defined in the Di config
func {{.receiver}}Show{{.name}}Action({{.in}}){{.out}} {
{{.return}}}

// Update{{.name}}Action function can take any parameters defined in the Di config
func {{.receiver}}Update{{.name}}Action({{.in}}){{.out}} {
{{.return}}}

// Destroy{{.name}}Action function can take any parameters defined in the Di config
func {{.receiver}}Destroy{{.name}}Action({{.in}}){{.out}} {
{{.return}}}
