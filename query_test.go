package fastapi

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

const colName = "inventory"

func TestBson(t *testing.T) {

	filterJSON := "{}"

	var filter bson.M
	err := bson.UnmarshalExtJSON([]byte(filterJSON), true, &filter)
	if err != nil {
		log.Println(err, filterJSON)
		return
	}

	log.Println(len(filter))

}

func TestFindWithBson(t *testing.T) {

	initTestDb()

	//param := `{"item":"mousepad"}`
	param := ""
	var result []bson.M
	option := &FindOptions{}
	_colName := colName
	option.CollectionName = &_colName
	option.Results = &result
	res := DbTestInstance.FindWithBson(param, *option)

	jsonbytes, _ := json.Marshal(res)

	fmt.Println(string(jsonbytes))

}

func TestFindOneWithBson(t *testing.T) {

	initTestDb()

	param := `{"item":"mousepad"}`

	var result bson.M
	opton := CreateFindOneOptions(colName)
	DbTestInstance.FindOneWithBson(param, &result, *opton)

	jsonbytes, _ := json.Marshal(&result)

	fmt.Println(string(jsonbytes))

}

func TestFindOne(t *testing.T) {

	initTestDb()

	var result Collection

	filter := bson.M{"name": "demo"}

	opton := CreateFindOneOptions("collection")
	DbTestInstance.FindOne(filter, &result, *opton)

	jsonbytes, _ := json.Marshal(&result)

	fmt.Println(opton, string(jsonbytes))

}

func TestFindId(t *testing.T) {

	initTestDb()

	var result Collection
	// oid, _ := primitive.ObjectIDFromHex("60ab75b9b6dcb68ed62efb11")
	// DbTestInstance.FindOne(bson.M{"_id": oid}, &result)

	DbTestInstance.FindById("60abb4c28b687f6d3d4febb6", &result, FindOneOptions{})

	jsonbytes, _ := json.Marshal(&result)

	fmt.Println(string(jsonbytes))

}

func TestQueryWithBson(t *testing.T) {

	initTestDb()

	//param := `{"item":"mousepad"}`
	param := `[{"$match":{"item":"mousepad"}}]`

	option := FindOptions{}
	option.SetCollectionName(colName)
	option.Results = &[]bson.M{}

	res := DbTestInstance.QueryWithBson(param, option)

	jsonbytes, _ := json.Marshal(res)

	fmt.Println(string(jsonbytes))

}
