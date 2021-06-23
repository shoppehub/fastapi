package types

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/shoppehub/fastapi/collection"
	"github.com/shoppehub/fastapi/crud"
	"go.mongodb.org/mongo-driver/bson"
)

func TestConvert(t *testing.T) {

	filterJSON := `{"_id":"ObjectId(\"60b4625808c3c857cf62b712\")"}`

	var filter bson.M
	err := bson.UnmarshalExtJSON([]byte(filterJSON), true, &filter)
	if err != nil {
		log.Println(err, filterJSON)
		return
	}

	log.Println(filter["_id"])

	// a := make(map[string]interface{})
	// a := 12.1

	// a := [...]string{"1", "2"}

	a := 12

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

	result, cerr := Convert(&crud.Resource{}, &m, col)
	if cerr != nil {
		fmt.Println(cerr)
	}

	bytes, _ := json.Marshal(result)

	fmt.Println(string(bytes))

}
