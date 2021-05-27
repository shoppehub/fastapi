package crud

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/shoppehub/fastapi/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDel(t *testing.T) {

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
		Filter:         []string{"name"},
		Inc: []Inc{{
			Key:   "age",
			Value: 1,
		}},
	})

	u.Profile.School.Name = "123"
	Id := primitive.NewObjectID()
	u.Id = &Id

	result, err = DbTestInstance.SaveOrUpdateOne(u, &UpdateOption{
		CollectionName: &tableName,
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

	var r3 []user
	curs, err := DbTestInstance.DB.Collection(tableName).Find(context.Background(), bson.M{})
	curs.All(context.Background(), &r3)

	log.Println(r3)

	// time.Sleep(time.Duration(2) * time.Second)

	DbTestInstance.DeleteById(tableName, *result.(*user).Id)

	var r2 *user
	DbTestInstance.DB.Collection(tableName).FindOne(context.Background(), bson.M{"_id": result.(*user).Id}).Decode(r2)

	if r2 != nil {
		t.Fatal("DeleteById error")
	}

	DbTestInstance.DeleteAny(&DeleteOption{
		CollectionName: tableName,
		Filter: []Filter{{
			Key:      "name",
			Operator: OP_IN.toString(),
			Value:    []string{"123456"},
		}},
	})

	var r4 []user
	curs, err = DbTestInstance.DB.Collection(tableName).Find(context.Background(), bson.M{})

	curs.All(context.Background(), &r4)
	if len(r4) > 0 {
		t.Fatal("DeleteAny error")
	}
	// jsonResult, err := json.Marshal(r2)
	// fmt.Println(string(jsonResult))

}
