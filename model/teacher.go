package model

import "time"

/**
学生
 */
type Teacher struct {
	Id int64 `sql:"id"`
	Name string `sql:"name"`
	Address string `sql:"address"`
	No string `sql:"no"`
	ClassId int64 `sql:"class_id"`
	Create time.Time `sql:"create_date"`
	IsReading bool `sql:"is_reading"`
	High float64 `sql:"high"`
	Weight float64 `sql:"weight"`
}
