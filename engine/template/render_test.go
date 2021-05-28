package template

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func TestRender(t *testing.T) {

	templ := `
	db.orders.aggregate([
		{ $match: { status: "A" } },
		{ $group: { _id: "$cust_id", total: { $sum: "$amount" } } },
		{ $sort: { total: -1 } }
	])
	
	{{ m := map("foo", "bar", "asd", 123)}}
{{ range k := m }}
    {{k}}: {{.}}
{{ end }}
	{{ pid := "123" }}
	{{if pid =="123"}}
		{{ pipeline("all","{\"name\":1}") }}
	{{end}}
	`

	template, err := views.Parse("demo", templ)
	if err != nil {
		logrus.Error(err)
	}

	var resp bytes.Buffer
	vars := make(jet.VarMap)

	// vars.SetFunc("base64", func(a jet.Arguments) reflect.Value {
	// 	// a.RequireNumOfArguments("base64", 1, 1)

	// 	buffer := bytes.NewBuffer(nil)
	// 	fmt.Fprint(buffer, a.Get(0))

	// 	return reflect.ValueOf(base64.URLEncoding.EncodeToString(buffer.Bytes()))
	// })

	vars.SetFunc("pipeline", func(a jet.Arguments) reflect.Value {
		// a.RequireNumOfArguments("base64", 1, 1)

		fmt.Println("out:", a.Get(0))
		jsonStr := a.Get(1).Interface().(string)
		var m bson.M
		bson.UnmarshalExtJSON([]byte(jsonStr), true, &m)
		fmt.Println("out:", m)

		return reflect.ValueOf("")
	})

	if err = template.Execute(&resp, vars, nil); err != nil {
		logrus.Error(err)

	}

	// fmt.Println(resp.String())

}
