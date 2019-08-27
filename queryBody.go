package go_orm

import "reflect"

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
