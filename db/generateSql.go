package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

//生产sql 及参数


/**
反射获取 属性不为空的字段
*/
func getNotValuePram(obj interface{}) (map[string]interface{},string)  {

	var ob reflect.Type
	var obValue reflect.Value

	//是否是指针
	if reflect.TypeOf(obj).Kind() == reflect.Ptr{
		ob = reflect.TypeOf(obj).Elem()
		obValue = reflect.ValueOf(obj).Elem()
	}else{
		ob = reflect.TypeOf(obj)
		obValue = reflect.ValueOf(obj)
	}


	fieldKV := make(map[string]interface{})

		for i:=0;i<ob.NumField();i++{
			//获取字段
			fType := ob.Field(i)
			if(fType.Type.Kind() == reflect.Struct){
				if(fType.Type.Name() != "Time"){
					continue
				}
			}
			f:=ob.Field(i).Name
			//获取字段value
			v := obValue.FieldByName(f)
			//获取是否有tag
			tag := ob.Field(i).Tag.Get("sql")
			if len(tag)>0 {
				f=tag
			}
			if isNotBlank(v){
				fieldKV[f] =conver(v,fType)
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
	para,tableName := getNotValuePram(o)
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
   // tags := getTag(o,sqlWhere...)
	value := make([]interface{},0,10)
	para,tableName := getNotValuePram(o)
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
    //tags := getTag(o,sqlwhere...)
	value := make([]interface{},0,10)
	para,tableName := getNotValuePram(o)
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
拼接查询
 */
func find(oType reflect.Type, findWhere map[string]interface{}, findFields ...string) (string, []interface{}, []string) {
	//felid := make([]string,0)
	value := make([]interface{},0,10)

	//获取 tags fields
	tags,fields := getTagAndFeild(oType)
	//结构体
	switch oType.Kind() {
	case reflect.Struct:
		//oType = reflect.Type(o).Kind()
	case reflect.Map:
	}

	sql:= "select "
	if len(findFields)>0 {
		for _,f := range findFields {
			sql += f +"  ,"
		}
	}else {
		for i,tag := range tags  {
			if tag=="NULL"{
				sql +=" "+fields[i]+","
			}else {
				sql +=" "+tag+","
			}
		}
	}
	sql=sql[:len(sql)-1]

	//todo  关联关系
    sql+= " from " + oType.Name()
	where := ""
    if(findFields!=nil || len(findWhere)>0){
		where = " where 1=1  "
		for k,v := range findWhere{
			where+="and "+k+"=? "
			value=append(value,v)
		}
	}


	returnFields:=findFields
	if len(findFields)==0{
		returnFields=fields
	}

	return sql+where,value,returnFields
}

type PageInfo struct {
	CurPage int64
	PageSize int64
	TotalRecord int64
	TotalPageNum int64
}


//分页 todo
/*func page(o interface{}, findWhere map[string]interface{}, findFields ...string) (string, []interface{}){
	limit (curPage-1)*pageSize,pageSize
}*/





/**
数据类型转换
 */
func conver(value reflect.Value,fType reflect.StructField) interface{}  {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Bool:
		if value.Bool(){
			return 1
		}
		return 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint()
	case reflect.Float32, reflect.Float64:
		return value.Float()
	case reflect.Struct:
		if fType.Type.Name()=="Time" {
			i := value.Interface().(time.Time)
			return i.Format("2006-01-02 15:04:05")
		}
		return nil
	default:
		return nil
	}
}


func mappingConver(columnTypes []*sql.ColumnType, results []interface{}) *[]interface{} {

	converResult := make([]interface{},0)

	for i := 0;i< len(columnTypes);i++ {
		re := string(*reflect.ValueOf(results[i]).Interface().(*sql.RawBytes))
		switch columnTypes[i].DatabaseTypeName(){
			case "VARCHAR","CHAR","TEXT"://字符串
				converResult= append(converResult,re)
			case "TIMESTAMP"://日期
				if len(re)==0 {
					converResult= append(converResult,time.Time{})
				}else{
					date,err :=time.Parse("2006-01-02 15:04:05",re)
					if err!=nil{
						panic(err)
					}
					converResult= append(converResult,date)
				}
			case "FLOAT","DOUBLE","DECIMAL"://浮点
				reFloat,_ := strconv.ParseFloat(re,64)
				converResult= append(converResult,reFloat)
			case "INT","LONG"://整数
				reInt,_ := strconv.ParseInt(re,10,64)
				converResult= append(converResult,reInt)
			}

		}


	return &converResult
}

/**
获取tag
 */
func getTag(o interface{},whereSql ...string) []string {
	tags := make([]string,0,10)
	ob := reflect.TypeOf(o)
	for _,field := range whereSql {

		sField,find := ob.FieldByName(field)
		if !find  {
			panic(field+ "field is not ")
		}
		tag := sField.Tag.Get("sql")
		if(len(tag)==0){
			tag=sField.Name
		}
		tags=append(tags,tag)
	}
	return tags
}

/**
获取tag feild
*/
func getTagAndFeild(t reflect.Type) ([]string,[]string) {
	tags := make([]string,0)
	fields := make([]string,0)
	ob := t
	for i:=0;i<t.NumField(); i++  {
		sField := ob.Field(i)

		tag := sField.Tag.Get("sql")
		if len(tag)==0 {
			tag="NULL"
		}
		tags=append(tags,tag)
		fields=append(fields,sField.Name)
	}
	return tags,fields
}
