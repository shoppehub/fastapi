package template

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/CloudyKit/jet/v6"
)

type stringLoader struct {
	templates map[string]string
}

// NewLoader returns an initialized loader serving the passed http.FileSystem.
func NewStringLoader() jet.Loader {
	templates := make(map[string]string)
	return &stringLoader{
		templates: templates,
	}
}

func (l *stringLoader) PutTemplateString(key string, template string) {
	l.templates[key] = template
}

// Open implements Loader.Open() on top of an http.FileSystem.
func (l *stringLoader) Open(name string) (io.ReadCloser, error) {
	return ioutil.NopCloser(strings.NewReader(l.templates[name])), nil
}

// Exists implements Loader.Exists() on top of an http.FileSystem by trying to open the file.
func (l *stringLoader) Exists(name string) bool {
	if l.templates[name] != "" {
		return true
	}
	return false
}
