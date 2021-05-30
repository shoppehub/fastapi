package template

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"github.com/shoppehub/fastapi/collection"
	"github.com/shoppehub/fastapi/crud"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var views *jet.Set

var fsLoader interface{}

func init() {
	fsLoader = NewStringLoader()
	httpfsLoader := fsLoader.(jet.Loader)
	views = jet.NewSet(
		httpfsLoader,
	)
}

// 根据名称进行匹配
func Render(resource *crud.Resource, collection collection.Collection, fnName string, c *gin.Context) (map[string]interface{}, error) {

	fun := collection.Functions[fnName]

	loader := fsLoader.(*stringLoader)
	if !loader.Exists(fnName) {
		loader.templates["/"+fnName] = fun.Template
	}

	view, err := views.GetTemplate(fnName)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var resp bytes.Buffer
	result := make(map[string]interface{})
	vars := newVars(resource, result)

	if fun.Params != nil {
		for _, param := range fun.Params {
			var body map[string]interface{}
			c.ShouldBindJSON(&body)
			if body[param.Name] != nil {
				vars.Set(param.Name, body[param.Name])
			}
		}
	}

	if err = view.Execute(&resp, vars, nil); err != nil {
		logrus.Error(err)
		return nil, err
	}

	fmt.Println(resp.String())

	fmt.Println(result)

	return result, nil
}

// 初始化模板上下文
func newVars(resource *crud.Resource, result map[string]interface{}) jet.VarMap {
	vars := make(jet.VarMap)

	vars.SetFunc("d", func(a jet.Arguments) reflect.Value {
		d := bson.D{}

		for i := 0; i < a.NumOfArguments(); i += 2 {

			d = append(d, bson.E{
				Key:   a.Get(i).String(),
				Value: a.Get(i + 1).Interface(),
			})
		}
		m := reflect.ValueOf(d)
		return m
	})

	vars.SetFunc("aggregate", func(a jet.Arguments) reflect.Value {
		p := mongo.Pipeline{}
		for i := 0; i < a.NumOfArguments(); i++ {
			p = append(p, a.Get(i).Interface().(bson.D))
		}
		m := reflect.ValueOf(p)

		return m
	})

	vars.SetFunc("m", func(a jet.Arguments) reflect.Value {
		d := bson.M{}

		for i := 0; i < a.NumOfArguments(); i += 2 {
			d[a.Get(i).String()] = a.Get(i + 1).Interface()
		}
		m := reflect.ValueOf(d)
		return m
	})

	vars.SetFunc("findOption", func(a jet.Arguments) reflect.Value {
		d := crud.FindOptions{}
		collectionName := a.Get(0).Interface().(string)
		d.CollectionName = &collectionName

		if a.NumOfArguments() > 2 {
			curPage := a.Get(1).Interface()
			if reflect.ValueOf(curPage).Kind() == reflect.Int {
				d.CurPage = int64(curPage.(int))
			} else {
				d.CurPage = curPage.(int64)
			}

			pageSize := a.Get(2).Interface()
			if reflect.ValueOf(pageSize).Kind() == reflect.Int {
				d.PageSize = int64(pageSize.(int))
			} else {
				d.PageSize = pageSize.(int64)
			}
		}
		m := reflect.ValueOf(d)
		return m
	})

	vars.SetFunc("query", func(a jet.Arguments) reflect.Value {

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
			r, _ := resource.QueryWithoutPaging(a.Get(0).Interface().(mongo.Pipeline), findOption)
			return reflect.ValueOf(r)
		} else {
			r := resource.Query(a.Get(0).Interface().(mongo.Pipeline), findOption)
			return reflect.ValueOf(r)
		}
	})

	vars.SetFunc("context", func(a jet.Arguments) reflect.Value {
		// r := *result

		result[a.Get(0).String()] = a.Get(1).Interface()
		return reflect.ValueOf(result)
	})
	return vars
}
