package template

import (
	"reflect"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"go.mongodb.org/mongo-driver/bson"
)

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

func sortFunc() jet.Func {
	return func(a jet.Arguments) reflect.Value {
		m := a.Get(0).Interface().([]bson.M)
		return reflect.ValueOf(Sort(a.Get(1).String(), m, a.Get(2).Bool()))
	}
}
