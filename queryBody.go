package go_orm

import (
	"reflect"
	"strings"
)

/**
查询体
 */
type QueryBody struct {
	IsSlice bool
	IsMap bool
	T interface{}
	Db *Db
	TableName string
	Fields []string
	Tags []string
	TagFieldMap map[string]string
	Ttype reflect.Type
	Tvalue reflect.Value
}
/**
构造queryBody
 */
func (q QueryBody)ConstQueryBody(o interface{},db Db) *QueryBody  {
	queryBody := &QueryBody{}
	queryBody.Db=&db
	queryBody.T=o
	queryBody.Ttype=reflect.TypeOf(queryBody.T)
	queryBody.Tvalue=reflect.ValueOf(queryBody.T)

	if queryBody.Ttype.Kind() == reflect.Ptr{
		queryBody.Ttype = queryBody.Ttype.Elem()
		queryBody.Tvalue = queryBody.Tvalue.Elem()
	}else{
		panic("请传递指针类型")
	}
	//判断是否是分片类型

	if queryBody.Ttype.Kind()==reflect.Slice{
		queryBody.IsSlice=true
	}

	if queryBody.IsSlice{
		queryBody.Ttype = queryBody.Tvalue.Type().Elem()
		tableNames := strings.Split(queryBody.Ttype.String(),".")
		if len(tableNames)>0 {
			queryBody.TableName = tableNames[len(tableNames)-1]
		}
		queryBody.Tvalue = reflect.New(queryBody.Ttype).Elem()
	}else{
		queryBody.TableName = queryBody.Ttype.Name()
	}

	return queryBody
}