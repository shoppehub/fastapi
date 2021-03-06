package service

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"time"

	"github.com/shoppehub/fastapi/base"
	"github.com/shoppehub/fastapi/collection"
	"github.com/shoppehub/fastapi/crud"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Save(resource *crud.Resource, collection collection.Collection, body CollectionBody) (map[string]interface{}, error) {

	collectionName := collection.GetCollectionName()
	hasDbResult := false

	if body.Filter != nil && len(body.Filter) > 0 {
		var dbResult map[string]interface{}

		if body.Filter[base.ID] != nil {

			if reflect.ValueOf(body.Filter[base.ID]).Type().Name() != "ObjectID" {
				oid, err := primitive.ObjectIDFromHex(body.Filter[base.ID].(string))
				if err != nil {
					return nil, err
				}
				body.Filter[base.ID] = oid
			}
			resource.FindOne(body.Filter, &dbResult, crud.FindOneOptions{
				CollectionName: collectionName,
			})
		} else {
			filterBytes, _ := json.Marshal(body.Filter)
			var filter primitive.M
			bson.UnmarshalExtJSON(filterBytes, true, &filter)
			resource.FindOne(filter, &dbResult, crud.FindOneOptions{
				CollectionName: collectionName,
			})
		}

		if dbResult != nil {
			body.Body[base.ID] = dbResult[base.ID]
			hasDbResult = true
			body.Filter = bson.M{base.ID: body.Body[base.ID]}
		}
	} else {
		if body.Body[base.ID] == nil {
			body.Filter = bson.M{base.ID: primitive.NewObjectID()}
		} else {
			body.Filter = bson.M{base.ID: body.Body[base.ID]}
		}
	}

	setElements := bson.D{}
	setOnInsertElements := bson.D{}
	unsetElements := bson.D{}
	now := time.Now()
	for _, field := range collection.Fields {
		value := body.Body[field.Name]

		if field.Type == "id" {
			value = crud.GenerateId(resource, field.IdKey, field.IdInitVal)
			setOnInsertElements = append(setOnInsertElements, bson.E{field.Name, value})
			continue
		}

		if value == "null" {
			unsetElements = append(unsetElements, bson.E{field.Name, ""})
			continue
		}

		if field.Value == "time.Now()" {
			value = now
		}
		if field.SetOnInsert {
			setOnInsertElements = append(setOnInsertElements, bson.E{field.Name, value})
			continue
		}

		if value == nil && field.DefaultValue != nil && !hasDbResult {
			value = field.DefaultValue
		}

		if value != nil {
			setElements = append(setElements, bson.E{
				field.Name, value,
			})
		}
	}

	update := bson.D{
		{"$set", setElements},
		{"$setOnInsert", setOnInsertElements},
	}

	if len(unsetElements) > 0 {
		update = append(update, bson.E{"$unset", unsetElements})
	}

	one, err := resource.DB.Collection(*collectionName).UpdateOne(
		context.Background(),
		body.Filter,
		update,
		&options.UpdateOptions{Upsert: &crud.Upsert})
	if err != nil {
		updateBytes, _ := json.Marshal(update)
		logrus.Error(err, string(updateBytes))
		return nil, err
	}
	result := make(map[string]interface{})
	if one.UpsertedID != nil {
		oid := one.UpsertedID.(primitive.ObjectID)
		resource.FindById(oid.Hex(), result, crud.FindOneOptions{
			CollectionName: collectionName,
		})
	} else {
		resource.FindOne(body.Filter, result, crud.FindOneOptions{
			CollectionName: collectionName,
		})
	}

	return result, nil
}
func Update(resource *crud.Resource, collection collection.Collection, body CollectionBody) (*mongo.UpdateResult, error) {

	collectionName := collection.GetCollectionName()
	upsertFalse := false
	if body.Filter == nil || len(body.Filter) == 0 {
		logrus.Error("Filter can not be nil!")
		return &mongo.UpdateResult{}, nil
	}

	setElements := bson.D{}
	unsetElements := bson.D{}
	now := time.Now()
	for _, field := range collection.Fields {
		value := body.Body[field.Name]
		if value == nil {
			continue
		}
		if field.Name == "int" && value == 0 {
			continue
		}
		if field.Type == "objectId" {
			if reflect.TypeOf(value).Name() == "string" {
				hex, err := primitive.ObjectIDFromHex(value.(string))
				if err != nil {
					return nil, nil
				}
				setElements = append(setElements, bson.E{Key: field.Name, Value: hex})
			} else {
				setElements = append(setElements, bson.E{Key: field.Name, Value: value})
			}
			continue
		}
		if value == "null" {
			unsetElements = append(unsetElements, bson.E{Key: field.Name, Value: ""})
			continue
		}

		if field.Value == "time.Now()" {
			value = now
		}
		if value != nil {
			setElements = append(setElements, bson.E{Key: field.Name, Value: value})
		}
	}
	update := bson.D{
		{"$set", setElements},
	}

	if len(unsetElements) > 0 {
		update = append(update, bson.E{Key: "$unset", Value: unsetElements})
	}
	updateResult, err := resource.DB.Collection(*collectionName).UpdateMany(
		context.Background(),
		body.Filter,
		update,
		&options.UpdateOptions{Upsert: &upsertFalse})
	if err != nil {
		updateBytes, _ := json.Marshal(update)
		logrus.Error(err, string(updateBytes))
		return nil, err
	}

	return updateResult, nil
}
