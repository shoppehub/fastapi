package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/shoppehub/fastapi/collection"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestConvert(t *testing.T) {

	// a := make(map[string]interface{})
	// a := 12.1

	// a := [...]string{"1", "2"}

	a := mongo.Pipeline{}

	fmt.Println(reflect.ValueOf(a).Type().Name())
	rval := reflect.ValueOf(a)

	fmt.Println(rval.Kind())

	jsonStr := `
	{
		"userName":"1",
		"age":"12",
		"time":"2020-08-01",
		"nick":["123","2"],
		"obj":{
			"gmt":"2020-08-02 12:12:12"
		}
	}
	`
	var m map[string]interface{}
	bb := []byte(jsonStr)
	jerr := json.Unmarshal(bb, &m)

	if jerr != nil {
		fmt.Println(jerr)
		return
	}

	col := collection.Collection{
		Fields: []collection.CollectionField{
			0: {
				Name: "userName",
			},
			1: {
				Name: "age",
				Type: "int",
			},
			2: {
				Name: "time",
				Type: "time",
			},
			3: {
				Name: "nick",
				Type: "string[]",
			},
			4: {
				Name: "obj",
				Type: "object",
				Fields: []collection.CollectionField{
					0: {
						Name: "gmt",
						Type: "time",
					},
				},
			},
		},
	}

	result, cerr := Convert(&m, col)
	if cerr != nil {
		fmt.Println(cerr)
	}

	bytes, _ := json.Marshal(result)

	fmt.Println(string(bytes))

}
