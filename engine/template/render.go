package template

import (
	"bytes"

	"github.com/CloudyKit/jet/v6"
	"github.com/shoppehub/fastapi/collection"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var views *jet.Set

type Match struct {
	Cache []bson.D
}

func (d *Match) put(operator string, key string, value interface{}) {

	switch operator {
	case "gt":
		d.Cache = append(d.Cache, bson.D{{key, bson.D{{"gt", value}}}})

	}

}

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
