package go_orm

import (
	"bytes"
	"database/sql"
	"fmt"
	"gopkg.in/guregu/null.v3"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
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
			/*if fType.Type.Kind() == reflect.Struct {
				if fType.Type.Name() != "Time" {
					continue
				}
			}*/
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
	case reflect.Struct:
		switch value.Type().String() {
		case NULL_Bool:
			return !value.Field(0).FieldByName(FIELD_Valid).Bool()
		case NULL_Float:
			return !value.Field(0).FieldByName(FIELD_Valid).Bool()
		case NULL_Int:
			return !value.Field(0).FieldByName(FIELD_Valid).Bool()
		case NULL_String:
			return !value.Field(0).FieldByName(FIELD_Valid).Bool()
		case NULL_Time:
			return !value.FieldByName(FIELD_Valid).Bool()
		}

	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func isNotBlank(value reflect.Value) bool {
	return !isBlank(value)
}
/**
拼接插入sql
 */
func insertSql(o interface{}, dbSetting *DbSetting) (string, []interface{}) {
	value := make([]interface{},0)
	para,tableName := getNotValuePram(o)
	sql := "insert into "+Format(tableName,dbSetting.tableFormat)+"(%s) values (%s)"
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
func getUpdateSql(o interface{},dbSetting *DbSetting,sqlWhere ...string) (string, []interface{})  {
	fieldMap := make(map[string]int)
	wherePara := make([]interface{},0,5)
   // tags := getTag(o,sqlWhere...)
	value := make([]interface{},0)
	para,tableName := getNotValuePram(o)
	sql := "update "+Format(tableName,dbSetting.tableFormat)+" set %s  %s"
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
func getDeleteSql(o interface{},dbSetting *DbSetting,sqlwhere ...string) (string, []interface{})  {
	fieldMap := make(map[string]int)
	wherePara := make([]interface{},0,5)
    //tags := getTag(o,sqlwhere...)
	value := make([]interface{},0)
	para,tableName := getNotValuePram(o)
	sql := "delete from "+Format(tableName,dbSetting.tableFormat)+" %s "

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
func find(oType reflect.Type, findWhere map[string]interface{}, dbSetting *DbSetting, findFields ...string) (string, []interface{}, []string) {
	//felid := make([]string,0)
	value := make([]interface{},0)

	var tags  []string
	var fields []string

	//获取 指定 tags fields
	if len(findFields)>0 {
		tags,fields = getTagByFeild(oType,findFields,dbSetting)
	}else{
		tags,fields = getTagAndFeild(oType,dbSetting)
	}

	//结构体
	switch oType.Kind() {
	case reflect.Struct:
		//oType = reflect.Type(o).Kind()
	case reflect.Map:
	}

	sql:= "select "
	for i,tag := range tags  {
		if tag=="NULL"{
			sql +=" "+fields[i]+","
		}else {
			sql +=" "+tag+","
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

const(
	TIME  = "time.Time"

	NULL_Int="null.Int"
	NULL_String="null.String"
	NULL_Time="null.Time"
	NULL_Float="null.Float"
	NULL_Bool="null.Bool"


	ZERO_Int="zero.Int"
	ZERO_String="zero.String"
	ZERO_Time="zero.Time"
	ZERO_Float="zero.Float"
	ZERO_Bool="zero.Bool"

	SQL_Int="sql.NullFloat"
	SQL_String="sql.NullString"
	SQL_Float="sql.NullFloat"
	SQL_Bool="sql.NullBool"

	FIELD_Int    ="Int64"
	FIELD_Float  ="Float64"
	FIELD_String ="String"
	FIELD_Bool   ="Bool"
	FIELD_Valid  ="Valid"
)

/**
数据类型转换  获取字段
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
		switch fType.Type.String() {
		case TIME:
			i := value.Interface().(time.Time)
			return i.Format("2006-01-02 15:04:05")
		case NULL_Bool:
			if value.Field(0).FieldByName(FIELD_Bool).Bool(){
				return 1
			}else {
				return 0
			}
		case NULL_Float:
			return value.Field(0).FieldByName(FIELD_Float).Float()
		case NULL_Int:
			return value.Field(0).FieldByName(FIELD_Int).Int()
		case NULL_String:
			return value.Field(0).FieldByName(FIELD_String).String()
		case NULL_Time:
			t := value.FieldByName("Time").Interface().(time.Time)
			return t.Format("2006-01-02 15:04:05")
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

func mappingConverMap(columnTypes []*sql.ColumnType, results *[]interface{}, tagField *map[string]string, tagType *map[string]reflect.StructField) *map[string]interface{} {
	fieldValueMap := make(map[string]interface{})
	for i := 0;i< len(columnTypes);i++ {
		result:=(*results)[i]
		re := string(*reflect.ValueOf(result).Interface().(*sql.RawBytes))

		switch columnTypes[i].DatabaseTypeName(){
		case "VARCHAR","CHAR","TEXT"://字符串
			tagFieldConver(columnTypes[i].Name(),&fieldValueMap,tagField,re,tagType,FIELD_String)
		case "TIMESTAMP"://日期
			/*if len(re)==0 {
				tagFieldConver(columnTypes[i].Name(),&fieldValueMap,tagField,time.Time{})
			}else{
				date,err :=time.Parse("2006-01-02 15:04:05",re)
				if err!=nil{
					panic(err)
				}*/
				tagFieldConver(columnTypes[i].Name(),&fieldValueMap,tagField,re,tagType,TIME)
			//}
		case "FLOAT","DOUBLE","DECIMAL"://浮点
			//reFloat,_ := strconv.ParseFloat(re,64)
			tagFieldConver(columnTypes[i].Name(),&fieldValueMap,tagField,re,tagType,FIELD_Float)
		case "INT","LONG"://整数
			//reInt,_ := strconv.ParseInt(re,10,64)
			tagFieldConver(columnTypes[i].Name(),&fieldValueMap,tagField,re,tagType,FIELD_Int)
		}

	}
	return &fieldValueMap
}

func tagFieldConver(tagName string,fieldValueMap *map[string]interface{},tagField *map[string]string,re string,tagType *map[string]reflect.StructField,tableType string) {
	var value interface{}
	ff := (*tagField)[tagName]

	if len(ff)>0{
		fieldType := (*tagType)[tagName]
		switch fieldType.Type.Kind() {
		case reflect.Struct:
			switch fieldType.Type.String() {
			case TIME:
				if len(re)==0 {
					value=time.Time{}
				}else{
					date,err :=time.Parse("2006-01-02 15:04:05",re)
					if err!=nil{
						panic(err)
					}
					value=date
				}
			case NULL_Time:
				if len(re)==0 {
					value = null.NewTime(time.Now(),false)
				}else {
					date,err :=time.Parse("2006-01-02 15:04:05",re)
					if err!=nil{
						panic(err)
					}
					value = null.NewTime(date,true)
				}
			case NULL_String:
				if len(re)==0 {
					value = null.NewString(re,false)
				}else {
					value = null.NewString(re,true)
				}
			case NULL_Float:
				if len(re)==0 {
					value = null.NewFloat(0,false)
				}else {
					converValue,_ := strconv.ParseFloat(re,64)
					value = null.NewFloat(converValue,true)
				}
			case NULL_Int:
				if len(re)==0 {
					value = null.NewInt(0,false)
				}else {
					converValue,_ := strconv.ParseInt(re,10,64)
					value = null.NewInt(converValue,true)
				}
			case NULL_Bool:
				if len(re)==0 {
					value = null.NewBool(false,false)
				}else {
					converValue,_ := strconv.ParseInt(re,10,64)
					if converValue==int64(1){
						value = null.NewBool(true,true)

					}else{
						value = null.NewBool(false,true)
					}
				}
			}
		default:
			switch tableType {
			case FIELD_Int:
				value,_ = strconv.ParseInt(re,10,64)
			case FIELD_Float:
				value,_ = strconv.ParseFloat(re,64)
			//case FIELD_Bool:
			case FIELD_String:
			   value=re
			}

		}

		(*fieldValueMap)[ff]=value
	}

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
func getTagAndFeild(t reflect.Type,dbSetting *DbSetting) ([]string,[]string) {
	tags := make([]string,0)
	fields := make([]string,0)
	ob := t
	for i:=0;i<t.NumField(); i++  {
		sField := ob.Field(i)

		tag := sField.Tag.Get("sql")
		if len(tag)==0 {
			tag=Format(sField.Name,dbSetting.fieldFormat)
		}
		tags=append(tags,tag)
		fields=append(fields,sField.Name)
	}
	return tags,fields
}
/**
获取tag feild map
 */
func getTagAndFeildMap(t reflect.Type, dbSetting *DbSetting) (*map[string]string, *map[string]string, *map[string]reflect.StructField, *map[string]reflect.StructField) {
	tagFeildMap := make(map[string]string)
	feildTagMap := make(map[string]string)
	tagTypeMap := make(map[string]reflect.StructField)
	feildTypeMap := make(map[string]reflect.StructField)
	ob := t
	for i:=0;i<t.NumField(); i++  {
		sField := ob.Field(i)

		tag := sField.Tag.Get("sql")
		if len(tag)==0 {
			tag=Format(sField.Name,dbSetting.fieldFormat)
		}
		tagFeildMap[tag]=sField.Name
		tagTypeMap[tag]=sField
		feildTagMap[sField.Name]=tag
		feildTypeMap[sField.Name]=sField

	}
	return &tagFeildMap,&feildTagMap,&tagTypeMap,&feildTypeMap
}

func getTagByFeild(t reflect.Type, findFeilds []string,dbSetting *DbSetting) ([]string,[]string) {
	tags := make([]string,0)
	fields := make([]string,0)
	ob := t
	for i:=0;i< len(findFeilds); i++  {
		sField, err := ob.FieldByName(findFeilds[i])
		if !err  {
			panic("请确认查询的结构体字段存在")
		}
		tag := sField.Tag.Get("sql")
		if len(tag)==0 {
			tag=Format(sField.Name,dbSetting.fieldFormat)
		}
		tags=append(tags,tag)
		fields=append(fields,sField.Name)
	}
	return tags,fields
}

//格式化
func Format(str string,formatType int) string {
	switch formatType {
	case HUMP_UNDERLINE:
		return UnderscoreName(str)
	}
	return str
}


// 驼峰式写法转为下划线写法
func UnderscoreName(name string) string {
	buffer := bytes.NewBufferString("")
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteRune('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}

	return buffer.String()
}

// 下划线写法转为驼峰写法
func CamelName(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}