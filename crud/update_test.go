package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/shoppehub/fastapi/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	base.BaseId `bson,inline`
	Name        string     `json:"name,omitempty"`
	Age         int64      `bson:"age" json:"age,omitempty" update:"inc"`
	Time        *time.Time `bson:"time,omitempty" json:"time,omitempty" update:"skip"`
	Profile     *profile   `json:"profile,omitempty"`
}

type profile struct {
	School school `json:"school,omitempty"`
	Page   string `json:"page,omitempty"`
}

type school struct {
	Name string `bson:"name" json:"name,omitempty"`
}

func TestStruct(t *testing.T) {

	initTestDb()
	defer closeTestDb()

	tableName := "test_user"
	DbTestInstance.DB.Collection(tableName).Drop(context.Background())
	tt := time.Now()
	id := primitive.NewObjectID()
	u := &user{
		Name: "123456",
		// Age:  1,
		Time: &tt,
		BaseId: base.BaseId{
			Id:        &id,
			CreatedAt: &tt,
		},
		Profile: &profile{
			School: school{
				Name: "杭州",
			},
			Page: "222",
		},
	}

	result, err := DbTestInstance.SaveOrUpdateOne(u, &UpdateOption{
		CollectionName: &tableName,
		Filter:         []string{0: "name"},
		Inc: []Inc{0: Inc{
			Key:   "age",
			Value: 1,
		}},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	if result.(*user).Time != nil {
		t.Fatal("Time must be zero")
		return
	}

	time.Sleep(time.Duration(2) * time.Second)
	r2, err := DbTestInstance.SaveOrUpdateOne(u, &UpdateOption{
		CollectionName: &tableName,
		Filter:         []string{0: "name"},
		Inc: []Inc{0: {
			Key:   "age",
			Value: 3,
		}},
	})
	fmt.Println(r2.(*user).Id, result.(*user).Id)

	if r2.(*user).Id.Hex() != result.(*user).Id.Hex() {
		t.Fatal("id must be same")
	}

	jsonResult, err := json.Marshal(r2)
	fmt.Println(string(jsonResult))
}
