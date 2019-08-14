package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"go_web_curd/db/conn"
	"reflect"
)



func Save(obj interface{}) int64  {
	sqlStr,para := insertSql(obj)
	return exe(sqlStr,para)
}

func Update(obj interface{},whereSql ...string) int64  {
	sqlStr,para := getUpdateSql(obj,whereSql...)
	return exe(sqlStr,para)
}


func Delete(obj interface{},whereSql ...string) int64  {
	sqlStr,para := getDeleteSql(obj,whereSql...)
	return exe(sqlStr,para)
}

func exe(sql string,para []interface{}) int64  {
	logrus.WithFields(logrus.Fields{}).Info(sql,para)
	stmt, err := conn.GetDB().Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	result,err := stmt.Exec(para...)
	if err !=nil {
		logrus.WithFields(logrus.Fields{}).Error(err)
		panic(err)
	}
	re,_ := result.RowsAffected()
	return re
}

/**
直接slq执行
 */
func nativeSqlUpdate(nativeSql string,parameters []interface{}) int64  {
	return exe(nativeSql,parameters)
}

func getPrt(o interface{}) interface{} {
	oType := reflect.TypeOf(o)
	if oType.Kind() != reflect.Ptr{
		return &o
	}else {
		return o
	}
}

func FindQuery(o interface{}, findWhere map[string]interface{}, findFields ...string) *[]interface{} {
	list := make([]interface{},0)
	oType := reflect.TypeOf(o)
	oValue := reflect.ValueOf(o)
	if oType.Kind() == reflect.Ptr{
		oType = oType.Elem()
		oValue = oValue.Elem()
	}else{
		panic("请传递指针类型")
	}
	sqlStr, para, fields := find(oType,findWhere,findFields...)

	fmt.Println(fields)

	logrus.WithFields(logrus.Fields{}).Info(sqlStr,para)
	stmt, err := conn.GetDB().Prepare(sqlStr)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows,err := stmt.Query(para...)
	if err !=nil {
		logrus.WithFields(logrus.Fields{}).Error(err)
		panic(err)
	}
	defer rows.Close()
	for rows.Next()  {
		dataTypes,err :=rows.ColumnTypes()
		if err != nil{
			panic(err)
		}
		values := make([]sql.RawBytes,len(dataTypes) )
		scans := make([]interface{}, len(dataTypes))

		for i := range values {
			scans[i] = &values[i]
		}
		err = rows.Scan(scans...)
		if err!=nil{
			//logrus.WithFields(logrus.Fields{}).Error(err)
			panic(err)
		}

		converResult := mappingConver(dataTypes,scans)
		bean := resultMapping(oValue,converResult,fields)
		list = append(list,bean)
	}
	return &list


}
/**
封装返回值
 */
func resultMapping(v reflect.Value, result *[]interface{}, fields []string) interface{} {
	//fmt.Println(*result)
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


