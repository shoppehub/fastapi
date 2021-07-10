package collection

import (
	"fmt"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMode(t *testing.T) {
	str := `{ "fields": ["1", "2"] }`
	var col bson.M
	err := bson.UnmarshalExtJSON([]byte(str), true, &col)

	if err != nil {
		fmt.Println(err, 1)
	}

	fields := col["fields"]

	fmt.Println(fields, reflect.TypeOf(fields).String())

}
