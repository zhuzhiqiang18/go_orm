package db

import (
	"fmt"
	"reflect"
)

/**
反射获取 字段
*/
func getPram(obj interface{}) (map[string]interface{},string)  {
	ob := reflect.TypeOf(obj)
	obValue := reflect.ValueOf(obj)
	fieldKV := make(map[string]interface{})
	if ob.Kind() == reflect.Struct {
		for i:=0;i<ob.NumField();i++{
			//获取字段
			f:=ob.Field(i).Name
			//获取字段value
			v := obValue.FieldByName(f)
			//获取是否有tag
			tag := ob.Field(i).Tag.Get("sql")
			if len(tag)>0 {
				f=tag
			}
			if isNotBlank(v){
				fieldKV[f] =conver(v)
			}
		}
	}
	return fieldKV,ob.Name()
}
/**
判断Value是否有值
*/
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func isNotBlank(value reflect.Value) bool {
	return !isBlank(value)
}
/**
拼接插入sql
 */
func insertSql(o interface{}) (string, []interface{}) {
	value := make([]interface{},0,10)
	para,tableName := getPram(o)
	sql := "insert into "+tableName+"(%s) values (%s)"
	sqlPrefix := ""
	sqlSuffix := ""
	for k,v := range para {
		sqlPrefix+= k+","
		sqlSuffix+= "?,"
		value = append(value,v)
	}

	return fmt.Sprintf(sql,sqlPrefix[:len(sqlPrefix)-1],sqlSuffix[:len(sqlSuffix)-1]) ,value
}
/**
拼接更新sql
 */
func getUpdateSql(o interface{},sqlWhere ...string) (string, []interface{})  {
	fieldMap := make(map[string]int)
	wherePara := make([]interface{},0,5)

	value := make([]interface{},0,10)
	para,tableName := getPram(o)
	sql := "update "+tableName+" set %s  %s"
	sqlPrefix := ""
	sqlSuffix := "where 1=1 "
	for _,s := range sqlWhere {
		sqlSuffix+= "and "+s+ " = ?  "
		if(fieldMap[s]<=0){
			fieldMap[s]=1
		}
	}

	for k,v := range para {
		if fieldMap[k] <=0  {
			sqlPrefix+= k+" = ? ,"
			value = append(value,v)
		}else {
			wherePara=append(wherePara,v)
		}

	}
	value=append(value,wherePara...)

	return fmt.Sprintf(sql,sqlPrefix[:len(sqlPrefix)-1],sqlSuffix[:len(sqlSuffix)]) ,value
}
/**
拼接删除sql
 */
func getDeleteSql(o interface{},sqlwhere ...string) (string, []interface{})  {
	fieldMap := make(map[string]int)
	wherePara := make([]interface{},0,5)

	value := make([]interface{},0,10)
	para,tableName := getPram(o)
	sql := "delete from "+tableName+" %s "

	sqlSuffix := "where 1=1 "
	for _,s := range sqlwhere {
		sqlSuffix+= "and "+s+ " = ?  "
		if(fieldMap[s]<=0){
			fieldMap[s]=1
		}

	}
	for k,v := range para {
		if fieldMap[k] <=0  {

		}else {
			wherePara=append(wherePara,v)
		}

	}
	value=wherePara

	return fmt.Sprintf(sql,sqlSuffix[:len(sqlSuffix)]) ,value
}
/**
数据类型转换
 */
func conver(value reflect.Value) interface{}  {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Bool:
		return value.Int()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint()
	case reflect.Float32, reflect.Float64:
		return value.Float()
	default:
		panic("只支持基本类型")
	}
}