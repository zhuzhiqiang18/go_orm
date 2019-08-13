package conn

import (
	"database/sql"
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

	log := logrus.New()
	//log.SetReportCaller(true)
	logrus.SetLevel(logrus.InfoLevel)
	log.WithFields(logrus.Fields{}).Info("DB COON ……")

}
func GetDB() *sql.DB {
	return db
}
