package dao

import (
	"go_web_curd/model"
	"go_web_curd/orm"
)

type StudentDao model.Student

func (dao StudentDao) Save(s model.Student) int64 {
	return orm.Save(s)
}

func (dao StudentDao) Update(s model.Student,whereSql ...string) int64 {
	return orm.Update(s,whereSql...)
}

func (dao StudentDao) Delete(s model.Student,whereSql ...string) int64 {
	return orm.Delete(s,whereSql...)
}