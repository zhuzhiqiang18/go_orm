package go_orm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhuzhiqiang18/go_orm/conn"
	"github.com/zhuzhiqiang18/go_orm/logger"
	"reflect"
	"strings"
)

type interfaceDb interface {
     Prepare(query string) (*sql.Stmt, error)
}


type Db struct {
	connDb *sql.DB
	connTx *sql.Tx
	isTX bool
	abstractDb interfaceDb
	setting *DbSetting
}




func getDb(connDb *sql.DB) *Db {
	var db = Db{}
	db.connDb=connDb
	db.connTx=nil
	db.isTX=false
	db.abstractDb=connDb
	db.setting=&DbSetting{db:db.connDb,fieldFormat:DEFAULE,tableFormat:DEFAULE}
	return &db
}

func  Open(User string, Password string, Host string, Port int64, DataBaseName string) (*Db, error) {
	db,err := conn.Open(User,Password,Host,Port,DataBaseName)
	return getDb(db),err
}

func (db *Db) DBSetting() *DbSetting {
	return db.setting
}

func  (db Db) Close() error {
	return db.connDb.Close()
}

func (db Db) Save(obj interface{}) (int64, int64) {

	sqlStr,para := insertSql(obj,db.setting)
	return db.exe(sqlStr,para)
}

func (db Db) Update(obj interface{}, whereSql ...string) int64 {
	sqlStr,para := getUpdateSql(obj,db.setting,whereSql...)

	affected,_:= db.exe(sqlStr,para)
	return affected
}


func (db Db) Delete(obj interface{}, whereSql ...string) int64 {
	sqlStr,para := getDeleteSql(obj,db.setting,whereSql...)
	affected,_:= db.exe(sqlStr,para)
	return affected
}



func (db Db) exe(sqlStr string, para []interface{}) (int64, int64) {
	logger.Debug(sqlStr,para)
	stmt, err := db.abstractDb.Prepare(sqlStr)
	if err!=nil {
		panic(err)
	}
	var result sql.Result
	defer stmt.Close()
	result,err = stmt.Exec(para...)
	if err!=nil {
		panic(err)
	}
	//改变行数
	var affected int64
	//最后插入的ID
	lastInsertId :=int64(0)
	affected,err = result.RowsAffected()
	if err!=nil {
		panic(err)
	}

	if strings.HasPrefix(sqlStr,"insert") {
		lastInsertId,err =result.LastInsertId()
		if err!=nil {
			panic(err)
		}
	}

	return affected,lastInsertId
}

/**
直接slq执行
 */
func (db Db) NativeSql(nativeSql string, parameters ...interface{}) int64 {
	affected,_:= db.exe(nativeSql,parameters)
	return affected
}

/**
供外部调用gql查询
 */
func (db Db) FindGql(gql *Gql) error {
	return db.queryGQL(gql)
}

/**
gql 查询
 */
func (db Db) queryGQL(gql *Gql) error {
	gql.QueryBody= QueryBody{}.ConstQueryBody(gql.t,db)
	return db.exeQuery(gql.QueryBody,gql.GetGql(db.setting),gql.para)
}


/**
供外部调用查询
*/
func (db Db) FindQuery(o interface{},sqlStr string, para ...interface{}) error {
	queryBody := QueryBody{}.ConstQueryBody(o,db)

	return db.exeQuery(queryBody,sqlStr,para)
}

/**
执行查询
 */
func (db Db) exeQuery(queryBody *QueryBody,sqlStr string, para []interface{}) error {
	logger.Debug(sqlStr,para)

	stmt, err := db.abstractDb.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows,err := stmt.Query(para...)
	if err !=nil {
		return err
	}

	tagField, _, tagType, _ := getTagAndFeildMap(queryBody.Ttype,db.setting)

	defer rows.Close()
	for rows.Next()  {
		dataTypes,err :=rows.ColumnTypes()
		if err != nil{
			return err
		}
		values := make([]sql.RawBytes,len(dataTypes) )
		scans := make([]interface{}, len(dataTypes))

		for i := range values {
			scans[i] = &values[i]
		}
		err = rows.Scan(scans...)
		if err!=nil{
			return err
		}


		converResult := mappingConverMap(dataTypes,&scans,tagField,tagType)
		bean := resultMappingFieldValueMap(queryBody.Tvalue,converResult)

		if !queryBody.IsSlice {
			return err
		}else {
			results :=indirect(reflect.ValueOf(queryBody.T))
			results.Set(reflect.Append(results,reflect.ValueOf(bean)))

		}
	}
	return err
}




func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}




func RowsMapping(rows sql.Row,oValue reflect.Value,fields []string,list []interface{})  {

}




func (db *Db)Begin() *Db {
	var newDb = Db{}
	newDb.connDb= db.connDb

	tx, err := db.connDb.Begin()
	if err != nil {
		panic(err)
	}

	newDb.connTx=tx
	newDb.isTX=true
	newDb.abstractDb=newDb.connTx
	newDb.setting=db.setting
	return &newDb
}


func (db *Db) Commit() error {
	if db.connTx ==nil && !db.isTX {
		panic("未获得事务")
	}
	return db.connTx.Commit()
}

func (db *Db) Rollback() error {
	if db.connTx ==nil && !db.isTX {
		panic("未开启事务")
	}
	return db.connTx.Rollback()
}