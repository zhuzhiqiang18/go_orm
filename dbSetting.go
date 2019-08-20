package go_orm

import (
	"database/sql"
)



const (
	DEFAULE =1
	HUMP=2
	HUMP_UNDERLINE=3
)

type DbSetting struct {
	db *sql.DB
	tableFormat int
	fieldFormat int
}

func (dbSetting *DbSetting) SetTableFormat(format int)  {
	dbSetting.tableFormat=format
}

func (dbSetting *DbSetting) SetFieldFormat(format int)  {
	dbSetting.fieldFormat=format
}

