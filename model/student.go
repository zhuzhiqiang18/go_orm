package model

/**
学生
 */
type Student struct {
	Id int `sql:"id"`
	Name string
	Address string
	No string
	ClassId int `sql:"class_id"`
}
