package wizard

import (
	"fmt"
	"sort"
	"strings"
)

func NewClassCreator() ClassCreator {
	return &ClassWizard{}
}

type ClassCreator interface {
	GetHelp() string
	SetHelpHeader(string)
	GetTemplate(map[string]string) string
	GetTemplateParams(map[string]string) map[string]string
	SetParameterInfos(map[string]ParameterInfo)
	SetOutParameterInfos(map[string]ParameterInfo)
	SetTemplates(map[string]string)
}

type ParameterInfo struct {
	Name       string
	Alias      string
	ImportPath string
}

type ClassWizard struct {
	helpHeader        string
	templates         map[string]string
	parameterInfos    map[string]ParameterInfo
	outParameterInfos map[string]ParameterInfo
}

func (c *ClassWizard) SetHelpHeader(s string) {
	c.helpHeader = s
}

func (c *ClassWizard) SetTemplates(t map[string]string) {
	c.templates = t
}

func (c *ClassWizard) SetParameterInfos(p map[string]ParameterInfo) {
	c.parameterInfos = p
}

func (c *ClassWizard) SetOutParameterInfos(p map[string]ParameterInfo) {
	c.outParameterInfos = p
}

func (c *ClassWizard) GetTemplate(flags map[string]string) string {
	for flag := range flags {
		if template, ok := c.templates[flag]; ok {
			return template
		}
	}

	if defaultTemplate, ok := c.templates[""]; ok {
		return defaultTemplate
	}
	return ""
}

func (c *ClassWizard) GetTemplateParams(flags map[string]string) map[string]string {
	return map[string]string{
		"imports": c.getImportsAsString(flags),
		"in":      c.getInputParamsAsString(flags),
		"out":     c.getOutputParamsAsString(flags),
		"return":  c.getReturnsAsString(flags),
	}
}

func (c *ClassWizard) getInputParamsAsString(flags map[string]string) string {
	if flag, ok := flags["in"]; ok {
		params := c.getParams(flag)
		return strings.Join(params, ", ")
	}
	return ""
}

func (c *ClassWizard) getOutputParamsAsString(flags map[string]string) string {
	if flag, ok := flags["out"]; ok {
		params := c.getOutParams(flag)
		if len(params) == 0 {
			return ""
		}
		outParams := strings.Join(params, ", ")
		if len(params) > 1 {
			return fmt.Sprintf(" (%s)", outParams)
		}
		return fmt.Sprintf(" %s", outParams)
	}
	return ""
}

func (c *ClassWizard) getReturnsAsString(flags map[string]string) string {
	if flag, ok := flags["out"]; ok {
		params := c.getOutAliases(flag)
		if len(params) == 0 {
			return ""
		}
		outParams := strings.Join(params, ", ")

		return fmt.Sprintf("     return %s\n", outParams)

	}
	return ""
}

func (c *ClassWizard) getParams(s string) []string {
	result := make([]string, 0)
	requestedParams := strings.Split(s, ",")
	for _, requestedParam := range requestedParams {
		paramResult, ok := c.getRequestedParam(requestedParam)
		if ok {
			result = append(result, paramResult.Alias+" "+paramResult.Name)
		}
		// we may want to handle not ok, return error?

	}

	return result
}

func (c *ClassWizard) getOutParams(s string) []string {
	result := make([]string, 0)
	requestedParams := strings.Split(s, ",")
	for _, requestedParam := range requestedParams {
		paramResult, ok := c.getOutRequestedParam(requestedParam)
		if ok {
			result = append(result, paramResult.Name)
		}
		// we may want to handle not ok, return error?

	}

	return result
}

func (c *ClassWizard) getOutAliases(s string) []string {
	result := make([]string, 0)
	requestedParams := strings.Split(s, ",")
	for _, requestedParam := range requestedParams {
		paramResult, ok := c.getOutRequestedParam(requestedParam)
		if ok {
			result = append(result, paramResult.Alias)
		}
		// we may want to handle not ok, return error?

	}

	return result
}

func (c *ClassWizard) getImportNames(s string, isOut bool) []string {
	result := make([]string, 0)
	requestedParams := strings.Split(s, ",")
	var paramResult *ParameterInfo
	var ok bool
	for _, requestedParam := range requestedParams {
		if isOut {
			paramResult, ok = c.getOutRequestedParam(requestedParam)
		} else {
			paramResult, ok = c.getRequestedParam(requestedParam)
		}

		if ok && paramResult.ImportPath != "" {
			result = append(result, paramResult.ImportPath)
		}
		// we may want to handle not ok, return error?
	}

	return result
}

func (c *ClassWizard) getRequestedParam(s string) (*ParameterInfo, bool) {
	params := c.getParameterInfos()

	if par, ok := params[s]; ok {
		return &par, true
	}
	return nil, false
}

func (c *ClassWizard) getOutRequestedParam(s string) (*ParameterInfo, bool) {
	params := c.getOutParameterInfos()

	if par, ok := params[s]; ok {
		return &par, true
	}
	return nil, false
}

func (c *ClassWizard) getImportsAsString(flags map[string]string) string {
	imports := c.getImports(flags)

	if len(imports) == 0 {
		return ""
	}
	sb := &strings.Builder{}
	c.strWriter(sb, "\nimport (\n")

	for _, imp := range imports {
		c.strWriter(sb, "     ", imp, "\n")

	}
	c.strWriter(sb, ")\n")

	return sb.String()
}

func (c *ClassWizard) getImports(flags map[string]string) []string {
	result := make([]string, 0)
	inAndOutParams := make([]string, 0)

	if flag, ok := flags["in"]; ok {
		inParams := c.getImportNames(flag, false)
		inAndOutParams = append(inAndOutParams, inParams...)
	}

	if flag, ok := flags["out"]; ok {
		inParams := c.getImportNames(flag, true)
		inAndOutParams = append(inAndOutParams, inParams...)
	}

	for _, inAndOutPar := range inAndOutParams {
		if !c.sliceContains(result, inAndOutPar) {
			result = append(result, inAndOutPar)
		}
	}

	return result
}

func (*ClassWizard) sliceContains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func (c *ClassWizard) getParameterInfos() map[string]ParameterInfo {
	return c.parameterInfos
}

func (c *ClassWizard) getOutParameterInfos() map[string]ParameterInfo {
	return c.outParameterInfos
}

func (c *ClassWizard) GetHelp() string {
	sb := &strings.Builder{}
	c.strWriter(sb, c.helpHeader)

	ins := c.getParameterInfos()
	outs := c.getOutParameterInfos()
	inParNames := make([]string, 0)
	outParNames := make([]string, 0)

	if len(c.templates) > 1 {
		c.strWriter(sb, "\nTemplate variations if set, otherwise it will be the default:\n")
	}

	for templateName := range c.templates {
		if templateName != "" {
			c.strWriter(sb, " -", templateName, "\n")
		}
	}

	if len(ins) > 0 {
		c.strWriter(sb, "\nOptional -in parameter values:\n")
	}

	for parName := range ins {
		inParNames = append(inParNames, parName)

	}

	sort.Strings(inParNames)
	for _, key := range inParNames {
		parInfo := ins[key]
		c.strWriter(sb, " * ", key, ": ", parInfo.Alias, " ", parInfo.Name, "\n")
	}

	if len(ins) > 0 {
		c.strWriter(sb, "Example -in=", strings.Join(inParNames, ","), "\n")
	}

	if len(outs) > 0 {
		c.strWriter(sb, "\nOptional -out parameter values:\n")
	}

	for parName := range outs {
		outParNames = append(outParNames, parName)
	}

	sort.Strings(outParNames)

	for _, key := range outParNames {
		parInfo := outs[key]
		c.strWriter(sb, " * ", key, ": ", parInfo.Name, "\n")
	}

	if len(outs) > 0 {
		c.strWriter(sb, "Example -out=", strings.Join(outParNames, ","), "\n")
	}

	return sb.String()
}

func (c *ClassWizard) strWriter(sb *strings.Builder, pars ...string) {
	for _, par := range pars {
		sb.WriteString(par)
	}
}
