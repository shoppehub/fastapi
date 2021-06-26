package template

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/shoppehub/fastapi/collection"
	"github.com/shoppehub/fastapi/crud"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func Render(resource *crud.Resource, collection collection.Collection, fnName string, body map[string]interface{}) (map[string]interface{}, error) {

	fun := collection.Functions[fnName]

	loader := fsLoader.(*stringLoader)
	if !loader.Exists(fnName) {
		loader.templates["/"+fnName] = fun.Template
	}

	// view, err := views.GetTemplate(fnName)
	view, err := views.Parse(fnName, fun.Template)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var resp bytes.Buffer
	result := make(map[string]interface{})
	vars := NewVars(resource, result)

	if fun.Params != nil {
		for _, param := range fun.Params {
			if body[param.Name] != nil {
				vars.Set(param.Name, body[param.Name])
			}
		}
	}

	if err = view.Execute(&resp, *vars, nil); err != nil {
		logrus.Error(err)
		return nil, err
	}

	// fmt.Println(resp.String())

	// fmt.Println(result)

	return result, nil
}

//排序  start
type MapsSort struct {
	Key     string
	MapList []bson.M
	Desc    bool
}

func (m *MapsSort) Len() int {
	return len(m.MapList)
}

func (m *MapsSort) Less(i, j int) bool {
	if m.Desc {
		return m.MapList[i][m.Key].(float64) > m.MapList[j][m.Key].(float64)
	} else {
		return m.MapList[i][m.Key].(float64) < m.MapList[j][m.Key].(float64)
	}
}

func (m *MapsSort) Swap(i, j int) {
	m.MapList[i], m.MapList[j] = m.MapList[j], m.MapList[i]
}

func Sort(key string, maps []bson.M, desc bool) []bson.M {
	mapsSort := MapsSort{}
	mapsSort.Key = key
	mapsSort.MapList = maps
	mapsSort.Desc = desc
	sort.Sort(&mapsSort)

	return mapsSort.MapList
}

// 初始化模板上下文
func NewVars(resource *crud.Resource, result map[string]interface{}) *jet.VarMap {
	vars := make(jet.VarMap)

	vars.SetFunc("string", func(a jet.Arguments) reflect.Value {

		if !a.Get(0).IsValid() {
			return reflect.ValueOf("")
		}

		name := a.Get(0).Type().Name()

		switch name {
		case "ObjectID":
			oid := a.Get(0).Interface().(primitive.ObjectID)
			return reflect.ValueOf(oid.Hex())
		case "int":
			return reflect.ValueOf(fmt.Sprint(a.Get(0).Interface().(int)))
		}

		return reflect.ValueOf(a.Get(0).Interface())
	})

	vars.SetFunc("append", func(a jet.Arguments) reflect.Value {
		name := a.Get(0).Type().Name()
		if name == "M" {
			m := a.Get(0).Interface().(bson.M)
			if m[a.Get(1).String()] != nil {
				val := append(m[a.Get(1).String()].([]bson.M), a.Get(2).Interface().(bson.M))
				m[a.Get(1).String()] = val
			} else {
				val := []bson.M{a.Get(2).Interface().(bson.M)}
				m[a.Get(1).String()] = val
			}
			return reflect.ValueOf(m)
		} else {
			m := a.Get(0).Interface().(map[string]interface{})
			if m[a.Get(1).String()] != nil {
				val := append(m[a.Get(1).String()].([]interface{}), a.Get(2).Interface())
				m[a.Get(1).String()] = val
			} else {
				val := []interface{}{a.Get(2).Interface()}
				m[a.Get(1).String()] = val
			}
			return reflect.ValueOf(m)
		}
	})

	vars.SetFunc("map", func(a jet.Arguments) reflect.Value {
		if a.NumOfArguments()%2 > 0 {
			return reflect.ValueOf(make(map[string]interface{}))
		}
		m := reflect.ValueOf(make(map[string]interface{}, a.NumOfArguments()/2))
		for i := 0; i < a.NumOfArguments(); i += 2 {

			m.SetMapIndex(a.Get(i), a.Get(i+1))
		}
		return m
	})

	vars.SetFunc("put", func(a jet.Arguments) reflect.Value {

		name := a.Get(0).Type().Name()

		if name == "M" {
			m := a.Get(0).Interface().(bson.M)
			m[a.Get(1).String()] = a.Get(2).Interface()
			return reflect.ValueOf(m)
		} else {
			m := a.Get(0).Interface().(map[string]interface{})
			m[a.Get(1).String()] = a.Get(2).Interface()
			return reflect.ValueOf(m)
		}
	})

	vars.SetFunc("context", func(a jet.Arguments) reflect.Value {
		result[a.Get(0).String()] = a.Get(1).Interface()
		return reflect.ValueOf(result)
	})

	initDataBase(resource, &vars)
	initMath(&vars)

	return &vars
}

func initMath(vars *jet.VarMap) {
	vars.SetFunc("sort", func(a jet.Arguments) reflect.Value {

		m := a.Get(0).Interface().([]bson.M)
		return reflect.ValueOf(Sort(a.Get(1).String(), m, a.Get(2).Bool()))
	})
	vars.SetFunc("ceil", func(a jet.Arguments) reflect.Value {
		value := a.Get(0).Interface()
		return reflect.ValueOf(int(math.Ceil(value.(float64))))
	})
	vars.SetFunc("floor", func(a jet.Arguments) reflect.Value {
		value := a.Get(0).Interface()
		return reflect.ValueOf(int(math.Floor(value.(float64))))
	})
}

func initDataBase(resource *crud.Resource, vars *jet.VarMap) {

	vars.SetFunc("m", func(a jet.Arguments) reflect.Value {
		d := bson.M{}

		for i := 0; i < a.NumOfArguments(); i += 2 {
			d[a.Get(i).String()] = a.Get(i + 1).Interface()
		}
		m := reflect.ValueOf(d)
		return m
	})

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

	vars.SetFunc("findOption", func(a jet.Arguments) reflect.Value {
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

	vars.SetFunc("save", func(a jet.Arguments) reflect.Value {
		collectionName := a.Get(0).Interface().(string)
		data := a.Get(1).Interface()
		result, _ := resource.SaveOrUpdateOne(data, &crud.UpdateOption{
			CollectionName: &collectionName,
		})
		return reflect.ValueOf(result)
	})

}
