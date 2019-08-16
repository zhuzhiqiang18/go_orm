package conn

import (
	"database/sql"
	"fmt"
)

var db *sql.DB
var dbConifg DbSourceConfig
/*func init() {
	dbConifg.dbSourceConfig("root","123456","127.0.0.1",3306,"go_test")
	var err error
	db,err = sql.Open("mysql",dbConifg.GetDns())
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

	logger := logrus.New()
	//logger.SetReportCaller(true)
	logrus.SetLevel(logrus.InfoLevel)
	logger.WithFields(logrus.Fields{}).Info("DB COON ……")

}*/

func  Open(User string, Password string, Host string, Port int64, DataBaseName string) (*sql.DB, error) {
	dbConifg.dbSourceConfig(User,Password,Host,Port,DataBaseName)
	var err error
	db,err = sql.Open("mysql",dbConifg.GetDns())

	err = db.Ping()

	//设置最大连接数 0不限制
	db.SetMaxOpenConns(0)
	//设置最大闲置连接数
	db.SetMaxIdleConns(10)

	return db,err
}

func GetDB() *sql.DB {
	return db
}

type DbSourceConfig struct {
	User string
	Password string
	Host string
	Port int64
	Table string
}

func (dbConfig *DbSourceConfig) dbSourceConfig(User string,Password string,Host string,Port int64,Table string )  {
	dbConfig.Host=Host
	dbConfig.Password=Password
	dbConfig.Port=Port
	dbConfig.User=User
	dbConfig.Table=Table
}

func GetDbSourceConfig() DbSourceConfig {
	return dbConifg
}

func (dbConfig DbSourceConfig) GetDns() string  {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",dbConfig.User,dbConfig.Password,dbConfig.Host,dbConfig.Port,dbConfig.Table)
}