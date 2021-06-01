package collection

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMode(t *testing.T) {
	str := `
	{"name":"chemball/chemical",
	"desc":"化学品模型","developer":[{"name":"焦哥","time":"2021-06-01","desc":"初始化"}],"extend":"common/baseId",
	"fields":[{"name":"_id","type":"objectId","title":"主键"},
	{"name":"createdAt","type":"time","title":"创建时间","setOnInsert":true,"value":"time.Now()"},
	{"name":"updatedAt","type":"time","title":"修改时间","value":"time.Now()"},{"name":"casId","type":"id","initVal":10000001,"title":"casId","desc":"化工球自定义的唯一id"},
	{"name":"casNo","type":"string","title":"业务类型","desc":"官方的cas号，比如 120-12-2"},
	{"name":"categoryIds","type":"objectId[]","title":"类目Id","desc":"一个化学品可以属于多个类目"},
	{"name":"status","type":"string","title":"化学品状态","desc":"draft/online/offline"},
	{"name":"baseInfo","type":"object","title":"基本信息","fields":[{"name":"activeSubstance","type":"string","title":"有效成分","desc":"比如葡萄籽提取物的有效成份是原花青素"}]},
	{"name":"pProperties","type":"object","title":"物化性质","fields":[{"name":"meltingPoint","type":"string","title":"熔点"}]},
	{"name":"exportInfo","type":"object","title":"出口信息","fields":[{"name":"exportRebate","type":"float","title":"出口退税率"}]},
	{"name":"cProperties","type":"object","title":"属性分类部分",
		"fields":[{"name":"appearance","type":"select","title":"状态",
			"selectOptions":[{"label":"固态","value":1},{"label":"液态","value":2}]},
			{"name":"issueDate","type":"time","title":"创建时间","desc":"记录第一次创建时间，格式是 2021-01-01 10:10:10"}]}]}
	`
	var col Collection
	json.Unmarshal([]byte(str), &col)

	fmt.Println(col)
}
