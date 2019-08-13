package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)
var db *sql.DB

func init() {
	dns := "root:123456@tcp(127.0.0.1:3306)/go_test"
	var err error
	db,err = sql.Open("mysql",dns)
	if err!=nil {
		logrus.WithFields(logrus.Fields{}).Error("数据库链接失败")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Error("PING ERR")
		panic(err)
	}
	//设置最大连接数 0不限制
	db.SetMaxOpenConns(0)
	//设置最大闲置连接数
	db.SetMaxIdleConns(10)
}

func GetDB() *sql.DB {
	return db
}

func Save(obj interface{}) int64  {
	sql,para := insertSql(obj)
	return exe(sql,para)
}

func Update(obj interface{},whereSql ...string) int64  {
	sql,para := getUpdateSql(obj,whereSql...)
	return exe(sql,para)
}


func Delete(obj interface{},whereSql ...string) int64  {
	sql,para := getDeleteSql(obj,whereSql...)
	return exe(sql,para)
}

func exe(sql string,para []interface{}) int64  {
	stmt, err := db.Prepare(sql)
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
