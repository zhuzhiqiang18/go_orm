package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"go_web_curd/db/conn"
)



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
	logrus.WithFields(logrus.Fields{}).Info(sql)
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


