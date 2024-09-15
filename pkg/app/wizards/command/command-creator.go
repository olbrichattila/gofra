package commandcreator

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/olbrichattila/gofra/pkg/app/args"
	"github.com/olbrichattila/gofra/pkg/app/view"
)

func New() CommandCreator {
	return &Creator{}
}

type CommandCreator interface {
	Construct(a args.CommandArger, v view.Viewer)
	Create(string, string, map[string]string) error
}

type Creator struct {
	a args.CommandArger
	v view.Viewer
}

func (c *Creator) Construct(a args.CommandArger, v view.Viewer) {
	c.a = a
	c.v = v
}

func (c *Creator) Create(templateName, savePath string, data map[string]string) error {
	commandName, err := c.a.Get(0)
	if err != nil {
		return fmt.Errorf("file name not provided")
	}

	err = c.createFolderIfNotExists(savePath)
	if err != nil {
		return fmt.Errorf("could not create directory")
	}

	fileName := fmt.Sprintf("%s/%s.go", savePath, commandName)
	if _, err := os.Stat(fileName); err == nil {
		return fmt.Errorf("file already exists")
	}

	mergedData := map[string]string{
		"name": c.cleanFileName(commandName),
	}

	for key, value := range data {
		mergedData[key] = value
	}

	tpaths := c.getTemplatePath()
	c.v.NewPath(tpaths...)
	err = c.v.RenderToFile(fileName, templateName, mergedData)
	if err != nil {
		return err
	}

	return nil
}

func (c *Creator) cleanFileName(fn string) string {
	sb := &strings.Builder{}
	words := strings.FieldsFunc(fn, func(r rune) bool {
		return r == ' ' || r == '-' || r == '_'
	})

	for _, word := range words {
		sb.WriteString(
			c.filterSpecialChars(
				c.lcFirst(word),
			),
		)
	}

	return sb.String()

}

func (*Creator) lcFirst(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

func (*Creator) filterSpecialChars(s string) string {

	re := regexp.MustCompile("[^a-zA-Z0-9]+")

	return re.ReplaceAllString(s, "")
}

func (*Creator) createFolderIfNotExists(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		return err
	} else if err != nil {
		return err
	}

	return nil
}

// GetCurrentFilePath gets the absolute path of the current file
func (*Creator) getTemplatePath() []string {
	_, b, _, _ := runtime.Caller(0) // returns full path of this file
	fparts := strings.Split(b, "/")
	fparts[0] = "/" + fparts[0]
	rootPath := fparts[0 : len(fparts)-4]
	rootPath = append(rootPath, "command-templates")

	return rootPath
}
