package tests

import (
	"fmt"
	"go_web_curd/db"
	"go_web_curd/model"
)

func TestSave(){
	var student model.Student
	student.Name="张三"
	student.Address="中国"
	student.No="123456"
	student.ClassId=1
	res:= db.Save(student)
	fmt.Println("改变行数",res)
}

func TestDelete(){
	var student model.Student
	student.ClassId=1
	res:= db.Delete(student,"class_id")
	fmt.Println("改变行数",res)
}


func TestUpdate(){
	var student model.Student
	student.Name="张三"
	student.No="00000000"
	res:= db.Update(student,"name")
	fmt.Println("改变行数",res)
}