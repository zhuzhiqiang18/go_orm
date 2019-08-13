package model

/**
学生
 */
type Student struct {
	Id int `sql:"id"`
	Name string `sql:"name"`
	Address string `sql:"address"`
	No string
	ClassId int `sql:"class_id"`
}
