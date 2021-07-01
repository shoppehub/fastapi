package template

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"github.com/shoppehub/fastapi/collection"
	"github.com/shoppehub/fastapi/crud"
	"github.com/shoppehub/sjet"
	"github.com/shoppehub/sjet/context"
	"github.com/shoppehub/sjet/engine"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var ENGINE *engine.TemplateEngine

func InitEngine(resource *crud.Resource) {
	if ENGINE == nil {
		ENGINE = sjet.CreateWithMem()
		InitAPIFunc(resource)
	}
}

// 根据名称进行匹配
func Render(collection collection.Collection, fnName string, c *gin.Context) (map[string]interface{}, error) {

	fun := collection.Functions[fnName]

	templateContext := context.InitTemplateContext(ENGINE, c)

	_, err := sjet.RenderMemTemplate(ENGINE, templateContext, c, fnName, fun.Template)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return *templateContext.Context, nil
}

func InitAPIFunc(resource *crud.Resource) {
	sjet.RegCustomFunc("query", func(c *gin.Context) jet.Func {
		return queryFunc(resource)
	})
	sjet.RegCustomFunc("findOption", func(c *gin.Context) jet.Func {
		return findOptionFunc(resource)
	})
	sjet.RegCustomFunc("save", func(c *gin.Context) jet.Func {
		return saveFunc(resource)
	})

	sjet.RegCustomFunc("sort", func(c *gin.Context) jet.Func {
		return sortFunc()
	})
	sjet.RegCustomFunc("getCollectionFields", func(c *gin.Context) jet.Func {
		return getCollectionFieldsFunc(resource)
	})
}

func findOptionFunc(resource *crud.Resource) jet.Func {

	return func(a jet.Arguments) reflect.Value {
		d := crud.FindOptions{}
		collectionName := a.Get(0).Interface().(string)
		d.CollectionName = &collectionName

		if a.NumOfArguments() > 2 {
			curPage := a.Get(1).Interface()
			k := a.Get(1).Kind()
			if k == reflect.Int {
				d.CurPage = int64(curPage.(int))
			} else if k == reflect.Int64 {
				d.CurPage = curPage.(int64)
			} else {
				v := curPage.(float64)
				d.CurPage = int64(v)
			}
			pageSize := a.Get(2).Interface()
			k2 := a.Get(2).Kind()
			if k2 == reflect.Int {
				d.PageSize = int64(pageSize.(int))
			} else if k2 == reflect.Int64 {
				d.PageSize = pageSize.(int64)
			} else {
				v := pageSize.(float64)
				d.PageSize = int64(v)
			}
		}
		m := reflect.ValueOf(d)
		return m

	}
}

func queryFunc(resource *crud.Resource) jet.Func {

	return func(a jet.Arguments) reflect.Value {

		typeName := reflect.ValueOf(a.Get(0).Interface()).Type().Name()
		findOption := a.Get(1).Interface().(crud.FindOptions)

		if typeName == "M" {
			if findOption.CurPage == 0 && findOption.PageSize == 0 {
				r := resource.FindWithoutPaging(a.Get(0).Interface().(bson.M), findOption)
				return reflect.ValueOf(r)
			} else {
				r := resource.Find(a.Get(0).Interface().(bson.M), findOption)
				return reflect.ValueOf(r)
			}
		}
		if findOption.CurPage == 0 && findOption.PageSize == 0 {
			r, _ := resource.QueryWithoutPaging(a.Get(0).Interface().([]bson.D), findOption)
			return reflect.ValueOf(r)
		} else {
			r := resource.Query(a.Get(0).Interface().([]bson.D), findOption)
			return reflect.ValueOf(r)
		}
	}
}

func saveFunc(resource *crud.Resource) jet.Func {

	return func(a jet.Arguments) reflect.Value {
		collectionName := a.Get(0).Interface().(string)
		data := a.Get(1).Interface()
		result, _ := resource.SaveOrUpdateOne(data, &crud.UpdateOption{
			CollectionName: &collectionName,
		})
		return reflect.ValueOf(result)
	}
}

func getCollectionFieldsFunc(resource *crud.Resource) jet.Func {

	return func(a jet.Arguments) reflect.Value {
		collectionName := a.Get(0).Interface().(string)

		fieldMap := collection.FindOneCollection(resource, collectionName)

		return reflect.ValueOf(fieldMap)
	}
}
