package crud

import (
	"context"
	"log"

	"github.com/shoppehub/conf"
)

var DbTestInstance *Resource

func initTestDb() {

	// docker exec -it --rm mongo mongo admin
	//	instance, err := NewDB("mongodb://127.0.0.1:10001/?authSource=admin", "admin")

	conf.Init("")

	instance, err := NewDB(conf.GetString("mongodb.url"), conf.GetString("mongodb.dbname"))

	if err != nil {
		log.Panicln(err)
	}
	DbTestInstance = instance
	DbTestInstance.DB.Collection("test_user").Drop(context.Background())

}

func closeTestDb() {
	DbTestInstance.Client.Disconnect(context.Background())
}
