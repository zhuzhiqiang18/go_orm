package model

import (
	"gopkg.in/guregu/null.v3"
)

/**
学生
 */
type Teacher struct {
	Id null.Int `sql:"id"`
	Name null.String `sql:"name"`
	Address null.String `sql:"address"`
	No null.String `sql:"no"`
	ClassId null.Int
	Create null.Time `sql:"create_date"`
	IsReading null.Bool `sql:"is_reading"`
	High null.Float `sql:"high"`
	Weight null.Float `sql:"weight"`
}
