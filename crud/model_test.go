package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestS(t *testing.T) {

	typ := reflect.StructOf([]reflect.StructField{
		{
			Name: "Height",
			Type: reflect.TypeOf(float64(0)),
			Tag:  `json:"height"`,
		},
		{
			Name: "Age",
			Type: reflect.TypeOf(int(0)),
			Tag:  `json:"age"`,
		},
	})
	v := reflect.New(typ).Elem()
	v.Field(0).SetFloat(0.4)
	v.Field(1).SetInt(2)
	s := v.Addr().Interface()

	w := new(bytes.Buffer)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}

	fmt.Printf("value: %+v\n", s)
	fmt.Printf("json:  %s", w.Bytes())

	r := bytes.NewReader([]byte(`{"height":1.5,"age":10}`))
	if err := json.NewDecoder(r).Decode(s); err != nil {
		panic(err)
	}
	fmt.Printf("value: %+v\n", s)
}

func TestMode(t *testing.T) {

	type DeUser struct {
		Name string
	}

	RegisterType(&DeUser{})

	ins := NewStruct("de_user")

	d := ins.(DeUser)

	d.Name = "123"

	fmt.Print(d)

	oid, _ := primitive.ObjectIDFromHex("60e5661075cd18c89c0791b5")
	m := primitive.M{
		"id": oid,
	}

	rs, err := json.Marshal(&m)
	fmt.Println(string(rs), err)
}
