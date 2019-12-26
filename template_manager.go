package scaffold

import (
	"fmt"
	"strings"

	"github.com/nortonlifelock/files"
	"github.com/pkg/errors"
)

// Template is used to hold a string format. It can either be populated by passing a path to a template, or populating it with a string argument
// it acts as a cleaner strings.Replace that can be chained
type Template struct {
	base string
}

// NewTemplate returns a template that is populated with the contents contained within the file pointed at by the arguments
func NewTemplate(basePath string, templatePath string) (template *Template, err error) {
	if len(templatePath) > 0 {
		var t string // string template from file
		filePath := fmt.Sprintf("%s/templates/%s.template", basePath, templatePath)

		t, err = files.GetStringFromFile(filePath)
		if err == nil {
			if len(t) > 0 {
				template = &Template{
					base: t,
				}
			} else {
				err = errors.New("NewTemplate - The length of the template is zero")
			}
		}
	} else {
		err = errors.New("NewTemplate - The length of the template path is zero")
	}
	return template, err
}

// NewTemplateEmpty returns a template that doesn't contain anything. Often used in coordination with UpdateBase to populate a template with a string instead
// of a file
func NewTemplateEmpty() (template *Template) {
	template = &Template{
		base: "",
	}
	return template
}

// Repl replaces some contents of the template and returns the same template so the method calls can be chained
func (t *Template) Repl(stringToReplace string, replacingValue string) *Template {
	if len(stringToReplace) > 0 {
		t.base = strings.Replace(t.base, stringToReplace, replacingValue, -1)
	}
	return t
}

// UpdateBase changes the contents of the template to the contents of the argument. Completely eradicates the contents of the old template
func (t *Template) UpdateBase(newBase string) {
	t.base = newBase
}

// Get returns the contents of the template
func (t *Template) Get() (template string) {
	return t.base
}
