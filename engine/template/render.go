package template

import (
	"bytes"

	"github.com/CloudyKit/jet/v6"
	"github.com/shoppehub/fastapi/collection"
	"github.com/sirupsen/logrus"
)

var views *jet.Set

func init() {
	httpfsLoader := NewStringLoader()
	views = jet.NewSet(
		httpfsLoader,
	)
}

// 根据名称进行匹配
func Render(templateName string, collection collection.Collection) (string, error) {

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
