package template

import (
	"bytes"

	"github.com/CloudyKit/jet/v6"
	"github.com/sirupsen/logrus"
)

var views *jet.Set

func init() {
	httpfsLoader := NewStringLoader()
	views = jet.NewSet(
		httpfsLoader,
	)
}

func Render(templateName string) (string, error) {

	view, err := views.GetTemplate(templateName)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	var resp bytes.Buffer
	vars := make(jet.VarMap)

	if err = view.Execute(&resp, vars, nil); err != nil {
		logrus.Error(err)
		return "", err
	}

	return "", nil
}
