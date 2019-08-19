package go_orm

import (
	"reflect"
)

/**
封装返回值
*/
func resultMapping(v reflect.Value, result *[]interface{}, fields []string) interface{} {
	for i:=0;i< len(fields);i++  {
		if v.FieldByName(fields[i]).Type().Kind()==reflect.Bool {
			if (*result)[i] == int64(1){
				v.FieldByName(fields[i]).Set(reflect.ValueOf(true))
			}else {
				v.FieldByName(fields[i]).Set(reflect.ValueOf(false))
			}
		}else{
			v.FieldByName(fields[i]).Set(reflect.ValueOf((*result)[i]))
		}

	}
	return v.Interface()
}

/**
封装查询结果
 */
func resultMappingFieldValueMap(v reflect.Value, result *map[string]interface{}) interface{} {
	for field,value := range *result{
		if v.FieldByName(field).Type().Kind()==reflect.Bool {
			if value == int64(1){
				v.FieldByName(field).Set(reflect.ValueOf(true))
			}else {
				v.FieldByName(field).Set(reflect.ValueOf(false))
			}
		}else{
			v.FieldByName(field).Set(reflect.ValueOf(value))
		}
	}
	return v.Interface()
}


