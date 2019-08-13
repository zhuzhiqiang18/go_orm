package dao

import (
	"go_web_curd/db"
	"go_web_curd/model"
)

type StudentDao model.Student

func (dao StudentDao) Save(s model.Student) int64 {
	return db.Save(s)
}

func (dao StudentDao) Update(s model.Student,whereSql ...string) int64 {
	return db.Update(s,whereSql...)
}

func (dao StudentDao) Delete(s model.Student,whereSql ...string) int64 {
	return db.Delete(s,whereSql...)
}