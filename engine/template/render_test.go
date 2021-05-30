package template

import (
	"testing"

	"github.com/shoppehub/fastapi/collection"
)

func TestRender(t *testing.T) {

	templ := `
		{{  name := "123"}}
		{{ filter := d( "$match",d("age",d("gt",1),"name",name) ) }}
		{{ limit := d( "$limit",1)  }}
		{{ categorys := aggregate(filter,limit) }}

		{{context("categorys",categorys)}}
	`

	col := collection.Collection{
		Functions: make(map[string]collection.Function),
	}
	fun := collection.Function{
		Template: templ,
	}

	col.Functions["demo"] = fun

	// result, err := Render(col, "demo", nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(result, "success")

}
