package go_orm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/zhuzhiqiang18/go_orm/conn"
	"github.com/zhuzhiqiang18/go_orm/log"
	"reflect"
	"strings"
)

//var connDb sql.DB

type Db struct {
	connDb *sql.DB
}

func getDb(connDb *sql.DB) *Db {
	return &Db{connDb:connDb}
}

func  Open(User string, Password string, Host string, Port int64, Table string) (*Db, error) {
	db,err := conn.Open(User,Password,Host,Port,Table)
	return getDb(db),err
}


func  (db Db) Close() error {
	return db.connDb.Close()
}

func (db Db) Save(obj interface{}) (int64, int64, error) {

	sqlStr,para := insertSql(obj)
	return db.exe(sqlStr,para)
}

func (db Db) Update(obj interface{}, whereSql ...string) (int64, error) {
	sqlStr,para := getUpdateSql(obj,whereSql...)

	affected,_,err:= db.exe(sqlStr,para)
	return affected,err
}


func (db Db) Delete(obj interface{}, whereSql ...string) (int64, error) {
	sqlStr,para := getDeleteSql(obj,whereSql...)

	affected,_,err:= db.exe(sqlStr,para)
	return affected,err
}



func (db Db) exe(sqlStr string, para []interface{}) (int64, int64, error) {
	log.Debug(sqlStr,para)
	stmt, err := db.connDb.Prepare(sqlStr)
	var result sql.Result
	defer stmt.Close()
	result,err = stmt.Exec(para...)
	//改变行数
	var affected int64
	//最后插入的ID
	lastInsertId :=int64(0)
	affected,err = result.RowsAffected()

	if strings.HasPrefix(sqlStr,"insert") {
		lastInsertId,err =result.LastInsertId()
	}

	return affected,lastInsertId,err
}

/**
直接slq执行
 */
func (db Db) NativeSql(nativeSql string, parameters ...interface{}) (int64, error) {
	affected,_,err:= db.exe(nativeSql,parameters)
	return affected,err
}



func (db Db) FindQuery(o interface{}, findWhere map[string]interface{}, findFields ...string) *[]interface{} {
	list := make([]interface{},0)
	oType := reflect.TypeOf(o)
	oValue := reflect.ValueOf(o)
	if oType.Kind() == reflect.Ptr{
		oType = oType.Elem()
		oValue = oValue.Elem()
	}else{
		panic("请传递指针类型")
	}
	//拼接sql
	sqlStr, para, fields := find(oType,findWhere,findFields...)

	log.Debug(sqlStr,para)
	stmt, err := db.connDb.Prepare(sqlStr)
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


